package identity

import (
	"goairmon/site/services/session"
	"goairmon/site/testhelpers"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

func TestPanicWhenNoSessionStore(t *testing.T) {
	defer _assertPanic(t)

	service := NewIdentityService(nil)

	ctx := &testhelpers.FakeContext{Values: make(map[string]interface{})}
	ctx.Set(service.Cfg.CtxKeyCookieStore, &sessions.CookieStore{})
	_, _ = service.getSessionFromCookie(ctx)
}

func TestPanicWhenNoCookieStore(t *testing.T) {
	defer _assertPanic(t)

	service := NewIdentityService(nil)

	ctx := &testhelpers.FakeContext{Values: make(map[string]interface{})}
	ctx.Set(service.Cfg.CtxKeySessionStore, &session.SessionStore{})
	_, _ = service.getSessionFromCookie(ctx)
}

func TestErrorWhenNoSession(t *testing.T) {
	service := NewIdentityService(nil)

	ctx := &testhelpers.FakeContext{Values: make(map[string]interface{})}
	ctx.Set(service.Cfg.CtxKeySessionStore, &session.SessionStore{})
	ctx.Set(service.Cfg.CtxKeyCookieStore, sessions.NewCookieStore([]byte("cookie-encryption")))

	_, err := service.getSessionFromCookie(ctx)
	if err == nil {
		t.Error("expected error")
	}
}

func TestErrorWhenSessionIdNotFoundInSessionStore(t *testing.T) {
	service := NewIdentityService(nil)

	ctx := &testhelpers.FakeContext{Values: make(map[string]interface{})}
	cfg := session.Config{ExpirationSecs: 60}
	sessionStore := session.NewSessionStore(cfg)
	cookieStore := &testhelpers.FakeCookieStore{Sessions: make(map[string]*sessions.Session)}
	cookieSession, _ := cookieStore.New(nil, service.Cfg.CookieStoreKeySession)

	ctx.Set(service.Cfg.CtxKeySessionStore, sessionStore)
	ctx.Set(service.Cfg.CtxKeyCookieStore, cookieStore)

	_, err := service.getSessionFromCookie(ctx)
	if err == nil {
		t.Error("expected error")
	}

	cookieSession.Values[service.Cfg.CookiesValueSessionKey] = "first_session_id"
	_ = cookieStore.Save(nil, nil, cookieSession)

	_, err = service.getSessionFromCookie(ctx)
	if err == nil {
		t.Error("expected error")
	}

	newSess, _ := sessionStore.NewOrExisting("first_session_id")

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

	ctx := &testhelpers.FakeContext{Values: make(map[string]interface{})}
	cfg := session.Config{ExpirationSecs: 60}
	sessionStore := session.NewSessionStore(cfg)
	cookieStore := &testhelpers.FakeCookieStore{Sessions: make(map[string]*sessions.Session)}
	cookieSession, _ := cookieStore.New(nil, service.Cfg.CookieStoreKeySession)
	cookieSession.Values[service.Cfg.CookiesValueSessionKey] = "first_session_id"
	_ = cookieStore.Save(nil, nil, cookieSession)
	newSess, _ := sessionStore.NewOrExisting("first_session_id")
	ctx.Set(service.Cfg.CtxKeySessionStore, sessionStore)
	ctx.Set(service.Cfg.CtxKeyCookieStore, cookieStore)

	_ = service.LoadCurrentSession()(testhelpers.EmptyHandler)(ctx)

	sess := ctx.Get(service.Cfg.CtxKeySession)
	if sess != newSess {
		t.Error("expected matching sessions")
	}
}

func TestRequireSession(t *testing.T) {
	service := NewIdentityService(nil)
	responseRecorder := httptest.NewRecorder()

	ctx := &testhelpers.FakeContext{
		Values:     make(map[string]interface{}),
		FakeWriter: responseRecorder,
	}

	nextHandler := func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "here's the next response")
	}

	_ = service.RequireSession(nil)(nextHandler)(ctx)

	if ctx.Response().Status != http.StatusUnauthorized {
		t.Error("should be returning unauthorized")
	}

	ctx.Set(service.Cfg.CtxKeySession, &session.Session{})

	_ = service.RequireSession(nil)(nextHandler)(ctx)

	if ctx.Response().Status != http.StatusOK {
		t.Error("should be returning 200")
	}
}

func _assertPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("The code did not panic")
	}
}
