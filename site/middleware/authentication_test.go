package middleware

import (
	"goairmon/site/session"
	"testing"

	"net/http/httptest"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

func TestPanicWhenNoSessionStore(t *testing.T) {
	defer _assertPanic(t)

	service := NewIdentityService(nil)

	ctx := &_fakeContext{values: make(map[string]interface{})}
	ctx.Set(service.cfg.CtxKeyCookieStore, &sessions.CookieStore{})
	_, _ = service.getSessionFromCookie(ctx)
}

func TestPanicWhenNoCookieStore(t *testing.T) {
	defer _assertPanic(t)

	service := NewIdentityService(nil)

	ctx := &_fakeContext{values: make(map[string]interface{})}
	ctx.Set(service.cfg.CtxKeySessionStore, &session.SessionStore{})
	_, _ = service.getSessionFromCookie(ctx)
}

func TestErrorWhenNoSession(t *testing.T) {
	service := NewIdentityService(nil)

	ctx := &_fakeContext{values: make(map[string]interface{})}
	ctx.Set(service.cfg.CtxKeySessionStore, &session.SessionStore{})
	ctx.Set(service.cfg.CtxKeyCookieStore, sessions.NewCookieStore([]byte("cookie-encryption")))

	_, err := service.getSessionFromCookie(ctx)
	if err == nil {
		t.Error("expected error")
	}
}

func TestErrorWhenSessionIdNotFoundInSessionStore(t *testing.T) {
	service := NewIdentityService(nil)

	ctx := &_fakeContext{values: make(map[string]interface{})}
	cfg := session.Config{ExpirationSecs: 60}
	sessionStore := session.NewSessionStore(cfg)
	cookieStore := &_fakeCookieStore{sessions: make(map[string]*sessions.Session)}
	cookieSession, _ := cookieStore.New(nil, service.cfg.CookieStoreKeySession)

	ctx.Set(service.cfg.CtxKeySessionStore, sessionStore)
	ctx.Set(service.cfg.CtxKeyCookieStore, cookieStore)

	_, err := service.getSessionFromCookie(ctx)
	if err == nil {
		t.Error("expected error")
	}

	cookieSession.Values[service.cfg.CookiesValueSessionKey] = "first_session_id"
	_ = cookieStore.Save(nil, nil, cookieSession)

	_, err = service.getSessionFromCookie(ctx)
	if err == nil {
		t.Error("expected error")
	}

	newSess, _ := sessionStore.NewOrExisting("first_session_id", "")

	sess, err := service.getSessionFromCookie(ctx)
	if err != nil {
		t.Error("unexpected error", err)
	}

	if sess != newSess {
		t.Error("expected matching sessions", sess, newSess)
	}
}

func TestIdentityMiddleware(t *testing.T) {
	service := NewIdentityService(nil)

	ctx := &_fakeContext{values: make(map[string]interface{})}
	cfg := session.Config{ExpirationSecs: 60}
	sessionStore := session.NewSessionStore(cfg)
	cookieStore := &_fakeCookieStore{sessions: make(map[string]*sessions.Session)}
	cookieSession, _ := cookieStore.New(nil, service.cfg.CookieStoreKeySession)
	cookieSession.Values[service.cfg.CookiesValueSessionKey] = "first_session_id"
	_ = cookieStore.Save(nil, nil, cookieSession)
	newSess, _ := sessionStore.NewOrExisting("first_session_id", "")
	ctx.Set(service.cfg.CtxKeySessionStore, sessionStore)
	ctx.Set(service.cfg.CtxKeyCookieStore, cookieStore)

	service.LoadCurrentSession(ctx)(nil)

	sess := ctx.Get(service.cfg.CtxKeySession)
	if sess != newSess {
		t.Error("expected matching sessions")
	}
}

func TestRequireSession(t *testing.T) {
	service := NewIdentityService(nil)
	responseRecorder := httptest.NewRecorder()

	ctx := &_fakeContext{
		values:     make(map[string]interface{}),
		fakeWriter: responseRecorder,
	}

	nextHandler := func(ctx echo.Context) error {
		return nil
	}

	next := service.RequireSession(ctx)(nextHandler)

	if next != nil {
		t.Error("the next handler should not be called")
	}

	ctx.Set(service.cfg.CtxKeySession, session.Session{})

	next = service.RequireSession(ctx)(nextHandler)

	if next == nil {
		t.Error("the next handler should be called")
	}
}

func _assertPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("The code did not panic")
	}
}
