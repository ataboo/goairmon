package middleware

import (
	"fmt"
	"goairmon/site/session"

	gorilla "github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

const CookieSessionId = "goairmon_session"
const CtxCookie = "cookie_store"
const CtxSessionStore = "session_store"
const CtxSession = "current_session"
const CookieValueSessionId = "session_id"

func IdentityMiddleware(c echo.Context) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		var session, _ = getSessionFromCookie(c)
		c.Set(CtxSession, session)

		//TODO: set user

		return next
	}
}

func getSessionFromCookie(c echo.Context) (*session.Session, error) {
	cookieStore, ok := c.Get(CtxCookie).(gorilla.Store)
	if !ok {
		panic("failed to find cookie store in context")
	}

	sessionStore, ok := c.Get(CtxSessionStore).(*session.SessionStore)
	if !ok {
		panic("failed to find session store in context")
	}

	cookieSession, err := cookieStore.Get(c.Request(), CookieSessionId)
	if err != nil || cookieSession.IsNew {
		return nil, fmt.Errorf("failed to get an existing cookie session")
	}

	sessionId, ok := cookieSession.Values[CookieValueSessionId].(string)
	if !ok {
		return nil, fmt.Errorf("cookie does not contain a session id")
	}

	return sessionStore.Find(sessionId)
}
