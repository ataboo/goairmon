package identity

import (
	"goairmon/business/services/session"
	"goairmon/site/testhelpers"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

func TestCreatesNewCookieSession(t *testing.T) {
	service := NewIdentityService(nil)

	ctx := &testhelpers.FakeContext{Values: make(map[string]interface{}), FakeWriter: httptest.NewRecorder()}
	service.sessionStore = &session.SessionStore{}
	service.cookieStore = sessions.NewCookieStore([]byte("cookie-encryption"))

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
	service.sessionStore = session.NewSessionStore(cfg)
	service.cookieStore = &testhelpers.FakeCookieStore{Sessions: make(map[string]*sessions.Session)}
	cookieSession, _ := service.cookieStore.New(nil, service.Cfg.CookieStoreKeySession)

	if cookieSess, ok := ctx.Get(CtxCookieSession).(*sessions.Session); ok && cookieSess != nil {
		t.Error("expected no cookie session")
	}

	if sess, ok := ctx.Get(CtxServerSession).(*session.Session); ok && sess != nil {
		t.Error("expected no session in context")
	}

	cookieSession.Values[CookiesValueSessionKey] = "first_session_id"
	_ = service.cookieStore.Save(nil, nil, cookieSession)

	_ = service.storeSessionsInContext(ctx)

	if ctx.Get(CtxServerSession).(*session.Session) != nil {
		t.Error("expected no session in context")
	}

	newSess, _ := service.sessionStore.NewOrExisting("first_session_id")

	_ = service.storeSessionsInContext(ctx)

	if ctx.Get(CtxServerSession).(*session.Session) != newSess {
		t.Error("expected session to be stored in context")
	}
}

func TestIdentityMiddleware(t *testing.T) {
	service := NewIdentityService(nil)

	ctx := &testhelpers.FakeContext{Values: make(map[string]interface{})}
	cfg := session.Config{ExpirationSecs: 60}
	service.sessionStore = session.NewSessionStore(cfg)
	service.cookieStore = &testhelpers.FakeCookieStore{Sessions: make(map[string]*sessions.Session)}
	cookieSession, _ := service.cookieStore.New(nil, service.Cfg.CookieStoreKeySession)
	cookieSession.Values[CookiesValueSessionKey] = "first_session_id"
	_ = service.cookieStore.Save(nil, nil, cookieSession)
	newSess, _ := service.sessionStore.NewOrExisting("first_session_id")

	_ = service.LoadCurrentSession()(testhelpers.EmptyHandler)(ctx)

	sess := ctx.Get(CtxServerSession)
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

	runFlag := false
	service.RequireSession(func(c echo.Context) error {
		runFlag = true
		return nil
	})(nextHandler)(ctx)

	if runFlag == false {
		t.Error("flag should be set")
	}

	ctx.Set(CtxServerSession, &session.Session{})

	_ = service.RequireSession(nil)(nextHandler)(ctx)

	if ctx.Response().Status != http.StatusOK {
		t.Error("should be returning 200")
	}
}

func TestRedirectWithoutSession(t *testing.T) {
	ctx := &testhelpers.FakeContext{
		FakeWriter: httptest.NewRecorder(),
		Values:     make(map[string]interface{}),
	}

	identity := NewIdentityService(&IdentityConfig{CookieStoreKeySession: "app-key", CookieStoreEncryptionKey: "encryption-key"})

	_ = identity.RedirectUsersWithoutSession("/destinationpath")(testhelpers.EmptyHandler)(ctx)

	if ctx.Response().Status != http.StatusSeeOther {
		t.Error("unexpected status", ctx.Response().Status)
	}

	if ctx.RedirectPath != "/destinationpath" {
		t.Error("unexpected redirect", ctx.RedirectPath)
	}

	ctx.Set(CtxServerSession, &session.Session{})

	_ = identity.RedirectUsersWithoutSession("/destinationpath2")(func(c echo.Context) error {
		return c.Redirect(200, "/notredirected")
	})(ctx)

	if ctx.Response().Status != 200 || ctx.RedirectPath != "/notredirected" {
		t.Error("expected next method call")
	}
}

func TestRedirectWithSession(t *testing.T) {
	ctx := &testhelpers.FakeContext{
		FakeWriter: httptest.NewRecorder(),
		Values:     make(map[string]interface{}),
	}

	ctx.Set(CtxServerSession, &session.Session{})

	identity := NewIdentityService(&IdentityConfig{CookieStoreKeySession: "app-key", CookieStoreEncryptionKey: "encryption-key"})
	_ = identity.RedirectUsersWithSession("/destinationpath")(testhelpers.EmptyHandler)(ctx)

	if ctx.Response().Status != http.StatusSeeOther {
		t.Error("unexpected status", ctx.Response().Status)
	}

	if ctx.RedirectPath != "/destinationpath" {
		t.Error("unexpected redirect", ctx.RedirectPath)
	}

	ctx.Set(CtxServerSession, nil)

	_ = identity.RedirectUsersWithSession("/destinationpath2")(func(c echo.Context) error {
		return c.Redirect(200, "/notredirected")
	})(ctx)

	if ctx.Response().Status != 200 || ctx.RedirectPath != "/notredirected" {
		t.Error("expected next method call")
	}
}

func TestStartNewSession(t *testing.T) {
	ctx := &testhelpers.FakeContext{
		FakeWriter: httptest.NewRecorder(),
		Values:     make(map[string]interface{}),
	}

	ctx.Set(CtxServerSession, &session.Session{})
	identity := NewIdentityService(&IdentityConfig{CookieStoreKeySession: "app-key", CookieStoreEncryptionKey: "encryption-key"})

	if _, err := identity.StartNewSession(ctx); err == nil {
		t.Error("expected already logged in error: ")
	}

	ctx.Set(CtxServerSession, nil)

	if _, err := identity.StartNewSession(ctx); err != nil {
		t.Error("unexpected error", err)
	}
}

func TestEndSession(t *testing.T) {
	ctx := &testhelpers.FakeContext{
		FakeWriter: httptest.NewRecorder(),
		Values:     make(map[string]interface{}),
	}

	identity := NewIdentityService(&IdentityConfig{CookieStoreKeySession: "app-key", CookieStoreEncryptionKey: "encryption-key"})

	if err := identity.EndSession(ctx); err == nil {
		t.Error("expected already logged out error: ")
	}

	ctx.Set(CtxServerSession, &session.Session{})

	if err := identity.EndSession(ctx); err != nil {
		t.Error("unexpected error", err)
	}
}

func _assertPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("The code did not panic")
	}
}

type _fakeServiceProvider struct {
	services map[string]interface{}
}

func (p *_fakeServiceProvider) Register(key string, service interface{}) {
	p.services[key] = service
}
