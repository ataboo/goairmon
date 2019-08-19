package identity

import (
	"fmt"
	"goairmon/site/helper"
	"goairmon/site/services/session"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	gorilla "github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

const (
	CookiesValueSessionKey = "session_id"
	CtxCookieSession       = helper.CtxCookieSession
	CtxKeyCookieStore      = helper.CtxKeyCookieStore
	CtxKeySession          = helper.CtxKeySession
	CtxKeySessionStore     = helper.CtxKeySessionStore
)

func NewIdentityService(cfg *IdentityConfig) *IdentityService {
	if cfg == nil {
		cfg = DefaultIdentityCfg()
	}

	cookieStore := sessions.NewCookieStore([]byte(cfg.CookieStoreEncryptionKey))
	sessionStore := session.NewSessionStore(session.Config{
		ExpirationSecs: 60 * 60 * 24 * 28,
		GCDelaySeconds: 60 * 60,
	})

	return &IdentityService{
		Cfg:          cfg,
		cookieStore:  cookieStore,
		sessionStore: sessionStore,
	}
}

func DefaultIdentityCfg() *IdentityConfig {
	return &IdentityConfig{
		CookieStoreKeySession:    "gowebapp_session",
		CookieStoreEncryptionKey: "cookie-secret",
	}
}

type ServiceProvider interface {
	Register(key string, service interface{})
}

type IdentityConfig struct {
	CookieStoreKeySession    string
	CookieStoreEncryptionKey string
}

type IdentityService struct {
	Cfg          *IdentityConfig
	sessionStore *session.SessionStore
	cookieStore  *sessions.CookieStore
}

func (i *IdentityService) LoadCurrentSession() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := i.storeSessionsInContext(c); err != nil {
				return err
			}

			return next(c)
		}
	}
}

func (i *IdentityService) RequireSession(onNoSession echo.HandlerFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, ok := c.Get(CtxKeySession).(*session.Session)
			if !ok || session == nil {
				if onNoSession == nil {
					return c.String(http.StatusUnauthorized, "Please login to access this route")
				}
				return onNoSession(c)
			}

			return next(c)
		}
	}
}

func (i *IdentityService) RedirectUsersWithoutSession(redirectPath string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, ok := c.Get(CtxKeySession).(*session.Session)
			if ok && session != nil {
				return next(c)
			}

			return c.Redirect(http.StatusSeeOther, redirectPath)
		}
	}
}

func (i *IdentityService) RedirectUsersWithSession(redirectPath string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, ok := c.Get(CtxKeySession).(*session.Session)
			if !ok || session == nil {
				return next(c)
			}

			return c.Redirect(http.StatusSeeOther, redirectPath)
		}
	}
}

func (i *IdentityService) StartNewSession(c echo.Context) error {
	session, ok := c.Get(CtxKeySession).(*session.Session)
	if ok && session != nil {
		return fmt.Errorf("already logged in")
	}

	sessionId := uuid.New()
	session, err := i.sessionStore.NewOrExisting(sessionId.String())
	if err != nil {
		return fmt.Errorf("failed to make new session: %s", err)
	}
	c.Set(CtxKeySession, session)

	cookieSession, err := i.cookieStore.Get(c.Request(), i.Cfg.CookieStoreKeySession)
	if err != nil {
		return fmt.Errorf("failed to add session cookie: %s", err)
	}

	cookieSession.Values[CookiesValueSessionKey] = sessionId.String()
	err = cookieSession.Save(c.Request(), c.Response().Writer)
	if err != nil {
		return fmt.Errorf("failed to save session cookie: %s", err)
	}

	return nil
}

func (i *IdentityService) EndSession(c echo.Context) error {
	session, ok := c.Get(CtxKeySession).(*session.Session)
	if !ok || session == nil {
		return fmt.Errorf("already logged out")
	}

	_ = i.sessionStore.Remove(session.Id)
	cookieSession, err := i.cookieStore.Get(c.Request(), i.Cfg.CookieStoreKeySession)
	if err != nil {
		return fmt.Errorf("failed to get session cookie: %s", err)
	}

	cookieSession.Values[CookiesValueSessionKey] = nil

	err = cookieSession.Save(c.Request(), c.Response().Writer)
	if err != nil {
		return fmt.Errorf("failed to save session cookie: %s", err)
	}

	return nil
}

func (i *IdentityService) storeSessionsInContext(c echo.Context) error {
	cookieStore, ok := c.Get(CtxKeyCookieStore).(gorilla.Store)
	if !ok {
		panic("failed to find cookie store in context")
	}

	sessionStore, ok := c.Get(CtxKeySessionStore).(*session.SessionStore)
	if !ok {
		panic("failed to find session store in context")
	}

	cookieSession, err := cookieStore.Get(c.Request(), i.Cfg.CookieStoreKeySession)
	if err != nil {
		return fmt.Errorf("failed to get cookie session")
	}
	if cookieSession.IsNew {
		if err = cookieSession.Save(c.Request(), c.Response().Writer); err != nil {
			return fmt.Errorf("failed to start new cookie session")
		}
	}

	c.Set(CtxCookieSession, cookieSession)

	sessionId, ok := cookieSession.Values[CookiesValueSessionKey].(string)
	if ok {
		sess, _ := sessionStore.Find(sessionId)
		c.Set(CtxKeySession, sess)
	}

	return nil
}

func (i *IdentityService) RegisterWithProvider(provider ServiceProvider) {
	provider.Register(CtxKeyCookieStore, i.cookieStore)
	provider.Register(CtxKeySessionStore, i.sessionStore)
}
