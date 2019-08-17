package middleware

import (
	"fmt"
	"goairmon/site/session"

	gorilla "github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

func NewIdentityService(cfg *IdentityConfig) *IdentityService {
	if cfg == nil {
		cfg = DefaultIdentityCfg()
	}

	return &IdentityService{
		cfg: cfg,
	}
}

func DefaultIdentityCfg() *IdentityConfig {
	return &IdentityConfig{
		CookieStoreKeySession:  "gowebapp_session",
		CookiesValueSessionKey: "session_id",
		CtxKeyCookieStore:      "cookie_store",
		CtxKeySessionStore:     "session_store",
		CtxKeySession:          "current_session",
	}
}

type IdentityConfig struct {
	CookieStoreKeySession  string
	CtxKeyCookieStore      string
	CtxKeySessionStore     string
	CtxKeySession          string
	CookiesValueSessionKey string
	OnNoSession            echo.HandlerFunc
}

type IdentityService struct {
	cfg *IdentityConfig
}

func (i *IdentityService) LoadCurrentSession(c echo.Context) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		var session, _ = i.getSessionFromCookie(c)
		c.Set(i.cfg.CtxKeySession, session)

		return next
	}
}

func (i *IdentityService) RequireSession(c echo.Context) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		session := c.Get(i.cfg.CtxKeySession)
		if session == nil {
			return i.cfg.OnNoSession
		}

		return next
	}
}

func (i *IdentityService) getSessionFromCookie(c echo.Context) (*session.Session, error) {
	cookieStore, ok := c.Get(i.cfg.CtxKeyCookieStore).(gorilla.Store)
	if !ok {
		panic("failed to find cookie store in context")
	}

	sessionStore, ok := c.Get(i.cfg.CtxKeySessionStore).(*session.SessionStore)
	if !ok {
		panic("failed to find session store in context")
	}

	cookieSession, err := cookieStore.Get(c.Request(), i.cfg.CookieStoreKeySession)
	if err != nil || cookieSession.IsNew {
		return nil, fmt.Errorf("failed to get an existing cookie session")
	}

	sessionId, ok := cookieSession.Values[i.cfg.CookiesValueSessionKey].(string)
	if !ok {
		return nil, fmt.Errorf("cookie does not contain a session id")
	}

	return sessionStore.Find(sessionId)
}
