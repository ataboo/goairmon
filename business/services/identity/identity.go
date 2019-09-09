package identity

import (
	"fmt"
	"goairmon/business/services/session"
	"goairmon/site/helper"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

const (
	CookiesValueSessionKey = "session_id"
	CtxCookieSession       = helper.CtxCookieSession
	CtxServerSession       = helper.CtxServerSession
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

type IdentityConfig struct {
	CookieStoreKeySession    string
	CookieStoreEncryptionKey string
}

type IdentityService struct {
	Cfg          *IdentityConfig
	sessionStore *session.SessionStore
	cookieStore  sessions.Store
}

func (i *IdentityService) LoadCurrentSession() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_ = i.storeSessionsInContext(c)

			return next(c)
		}
	}
}

func (i *IdentityService) RequireSession(onNoSession echo.HandlerFunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session, ok := c.Get(CtxServerSession).(*session.Session)
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
			session, ok := c.Get(CtxServerSession).(*session.Session)
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
			session, ok := c.Get(CtxServerSession).(*session.Session)
			if !ok || session == nil {
				return next(c)
			}

			return c.Redirect(http.StatusSeeOther, redirectPath)
		}
	}
}

func (i *IdentityService) StartNewSession(c echo.Context) (*session.Session, error) {
	session, ok := c.Get(CtxServerSession).(*session.Session)
	if ok && session != nil {
		return nil, fmt.Errorf("already logged in")
	}

	sessionID := uuid.New()
	session, err := i.sessionStore.NewOrExisting(sessionID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to make new session: %s", err)
	}

	err = i.setSessionIdInCookies(c, sessionID)
	if err != nil {
		i.sessionStore.Remove(sessionID.String())
		return nil, err
	}

	c.Set(CtxServerSession, session)

	return session, nil
}

func (i *IdentityService) setSessionIdInCookies(c echo.Context, sessionID uuid.UUID) error {
	cookieSession, err := i.cookieStore.Get(c.Request(), i.Cfg.CookieStoreKeySession)
	if err != nil {
		return fmt.Errorf("failed to add session cookie: %s", err)
	}

	cookieSession.Values[CookiesValueSessionKey] = sessionID.String()
	err = cookieSession.Save(c.Request(), c.Response().Writer)
	if err != nil {
		return fmt.Errorf("failed to save session cookie: %s", err)
	}

	return nil
}

func (i *IdentityService) EndSession(c echo.Context) error {
	session, ok := c.Get(CtxServerSession).(*session.Session)
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
	cookieSession, err := i.cookieStore.Get(c.Request(), i.Cfg.CookieStoreKeySession)
	if err != nil {
		i.cookieStore.New(c.Request(), i.Cfg.CookieStoreKeySession)
	}
	if cookieSession.IsNew {
		if err = cookieSession.Save(c.Request(), c.Response().Writer); err != nil {
			return fmt.Errorf("failed to start new cookie session")
		}
	}

	c.Set(CtxCookieSession, cookieSession)

	sessionID, ok := cookieSession.Values[CookiesValueSessionKey].(string)
	if ok {
		sess, _ := i.sessionStore.Find(sessionID)
		c.Set(CtxServerSession, sess)
	}

	return nil
}
