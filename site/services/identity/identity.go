package identity

import (
	"fmt"
	"goairmon/site/services/session"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	gorilla "github.com/gorilla/sessions"
	"github.com/labstack/echo"
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
		CookiesValueSessionKey:   "session_id",
		CtxKeyCookieStore:        "cookie_store",
		CtxKeySessionStore:       "session_store",
		CtxKeySession:            "current_session",
		CookieStoreEncryptionKey: "cookie-secret",
	}
}

type ServiceProvider interface {
	Register(key string, service interface{})
}

type IdentityConfig struct {
	CookieStoreKeySession    string
	CtxKeyCookieStore        string
	CtxKeySessionStore       string
	CtxKeySession            string
	CookiesValueSessionKey   string
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
			var session, _ = i.getSessionFromCookie(c)
			c.Set(i.Cfg.CtxKeySession, session)

			return next(c)
		}
	}
}

func (i *IdentityService) RequireSession(onNoSession echo.HandlerFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, ok := c.Get(i.Cfg.CtxKeySession).(*session.Session)
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

func (i *IdentityService) StartNewSession(c echo.Context) error {
	session := c.Get(i.Cfg.CtxKeySession).(*session.Session)
	if session != nil {
		return fmt.Errorf("already logged in")
	}

	sessionId := uuid.New()
	session, err := i.sessionStore.NewOrExisting(sessionId.String())
	if err != nil {
		return fmt.Errorf("failed to make new session: %s", err)
	}
	c.Set(i.Cfg.CtxKeySession, session)

	cookieSession, err := i.cookieStore.Get(c.Request(), i.Cfg.CookieStoreKeySession)
	if err != nil {
		return fmt.Errorf("failed to add session cookie: %s", err)
	}

	cookieSession.Values[i.Cfg.CookiesValueSessionKey] = sessionId.String()
	err = cookieSession.Save(c.Request(), c.Response().Writer)
	if err != nil {
		return fmt.Errorf("failed to save session cookie: %s", err)
	}

	return nil
}

func (i *IdentityService) EndSession(c echo.Context) error {
	session, ok := c.Get(i.Cfg.CtxKeySession).(*session.Session)
	if !ok || session == nil {
		return fmt.Errorf("already logged out")
	}

	_ = i.sessionStore.Remove(session.Id)
	cookieSession, err := i.cookieStore.Get(c.Request(), i.Cfg.CookieStoreKeySession)
	if err != nil {
		return fmt.Errorf("failed to get session cookie: %s", err)
	}

	cookieSession.Values[i.Cfg.CookiesValueSessionKey] = nil

	err = cookieSession.Save(c.Request(), c.Response().Writer)
	if err != nil {
		return fmt.Errorf("failed to save session cookie: %s", err)
	}

	return nil
}

func (i *IdentityService) getSessionFromCookie(c echo.Context) (*session.Session, error) {
	cookieStore, ok := c.Get(i.Cfg.CtxKeyCookieStore).(gorilla.Store)
	if !ok {
		panic("failed to find cookie store in context")
	}

	sessionStore, ok := c.Get(i.Cfg.CtxKeySessionStore).(*session.SessionStore)
	if !ok {
		panic("failed to find session store in context")
	}

	cookieSession, err := cookieStore.Get(c.Request(), i.Cfg.CookieStoreKeySession)
	if err != nil || cookieSession.IsNew {
		return nil, fmt.Errorf("failed to get an existing cookie session")
	}

	sessionId, ok := cookieSession.Values[i.Cfg.CookiesValueSessionKey].(string)
	if !ok {
		return nil, fmt.Errorf("cookie does not contain a session id")
	}

	return sessionStore.Find(sessionId)
}

func (i *IdentityService) RegisterWithProvider(provider ServiceProvider) {
	provider.Register(i.Cfg.CtxKeyCookieStore, i.cookieStore)
	provider.Register(i.Cfg.CtxKeySessionStore, i.sessionStore)
}
