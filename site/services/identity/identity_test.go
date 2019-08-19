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
	ctx.Set(CtxKeyCookieStore, &sessions.CookieStore{})
	_ = service.storeSessionsInContext(ctx)
}

func TestPanicWhenNoCookieStore(t *testing.T) {
	defer _assertPanic(t)

	service := NewIdentityService(nil)

	ctx := &testhelpers.FakeContext{Values: make(map[string]interface{})}
	ctx.Set(CtxKeySessionStore, &session.SessionStore{})
	_ = service.storeSessionsInContext(ctx)
}

func TestCreatesNewCookieSession(t *testing.T) {
	service := NewIdentityService(nil)

	ctx := &testhelpers.FakeContext{Values: make(map[string]interface{}), FakeWriter: httptest.NewRecorder()}
	ctx.Set(CtxKeySessionStore, &session.SessionStore{})
	ctx.Set(CtxKeyCookieStore, sessions.NewCookieStore([]byte("cookie-encryption")))

	err := service.storeSessionsInContext(ctx)
	if err != nil {
		t.Error("unexpected error")
	}

	if ctx.Get(CtxCookieSession).(*sessions.Session) == nil {
		t.Error("expected cookie session to be saved")
	}
}

func TestErrorWhenSessionIdNotFoundInSessionStore(t *testing.T) {
	service := NewIdentityService(nil)

	ctx := &testhelpers.FakeContext{Values: make(map[string]interface{})}
	cfg := session.Config{ExpirationSecs: 60}
	sessionStore := session.NewSessionStore(cfg)
	cookieStore := &testhelpers.FakeCookieStore{Sessions: make(map[string]*sessions.Session)}
	cookieSession, _ := cookieStore.New(nil, service.Cfg.CookieStoreKeySession)

	ctx.Set(CtxKeySessionStore, sessionStore)
	ctx.Set(CtxKeyCookieStore, cookieStore)

	if cookieSess, ok := ctx.Get(CtxCookieSession).(*sessions.Session); ok && cookieSess != nil {
		t.Error("expected no cookie session")
	}

	if sess, ok := ctx.Get(CtxKeySession).(*session.Session); ok && sess != nil {
		t.Error("expected no session in context")
	}

	cookieSession.Values[CookiesValueSessionKey] = "first_session_id"
	_ = cookieStore.Save(nil, nil, cookieSession)

	_ = service.storeSessionsInContext(ctx)

	if ctx.Get(CtxKeySession).(*session.Session) != nil {
		t.Error("expected no session in context")
	}

	newSess, _ := sessionStore.NewOrExisting("first_session_id")

	_ = service.storeSessionsInContext(ctx)

	if ctx.Get(CtxKeySession).(*session.Session) != newSess {
		t.Error("expected session to be stored in context")
	}
}

func TestIdentityMiddleware(t *testing.T) {
	service := NewIdentityService(nil)

	ctx := &testhelpers.FakeContext{Values: make(map[string]interface{})}
	cfg := session.Config{ExpirationSecs: 60}
	sessionStore := session.NewSessionStore(cfg)
	cookieStore := &testhelpers.FakeCookieStore{Sessions: make(map[string]*sessions.Session)}
	cookieSession, _ := cookieStore.New(nil, service.Cfg.CookieStoreKeySession)
	cookieSession.Values[CookiesValueSessionKey] = "first_session_id"
	_ = cookieStore.Save(nil, nil, cookieSession)
	newSess, _ := sessionStore.NewOrExisting("first_session_id")
	ctx.Set(CtxKeySessionStore, sessionStore)
	ctx.Set(CtxKeyCookieStore, cookieStore)

	_ = service.LoadCurrentSession()(testhelpers.EmptyHandler)(ctx)

	sess := ctx.Get(CtxKeySession)
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

	ctx.Set(CtxKeySession, &session.Session{})

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
