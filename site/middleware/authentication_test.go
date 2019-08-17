package middleware

import (
	"goairmon/site/session"
	"testing"

	"github.com/gorilla/sessions"
)

func TestPanicWhenNoSessionStore(t *testing.T) {
	defer _assertPanic(t)

	ctx := &_fakeContext{values: make(map[string]interface{})}
	ctx.Set(CtxCookie, &sessions.CookieStore{})
	_, _ = getSessionFromCookie(ctx)
}

func TestPanicWhenNoCookieStore(t *testing.T) {
	defer _assertPanic(t)

	ctx := &_fakeContext{values: make(map[string]interface{})}
	ctx.Set(CtxSessionStore, &session.SessionStore{})
	_, _ = getSessionFromCookie(ctx)
}

func TestErrorWhenNoSession(t *testing.T) {
	ctx := &_fakeContext{values: make(map[string]interface{})}
	ctx.Set(CtxSessionStore, &session.SessionStore{})
	ctx.Set(CtxCookie, sessions.NewCookieStore([]byte("cookie-encryption")))

	_, err := getSessionFromCookie(ctx)
	if err == nil {
		t.Error("expected error")
	}
}

func TestErrorWhenSessionIdNotFoundInSessionStore(t *testing.T) {
	ctx := &_fakeContext{values: make(map[string]interface{})}
	cfg := session.Config{}
	sessionStore := session.NewSessionStore(cfg)
	cookieStore := &_fakeCookieStore{sessions: make(map[string]*sessions.Session)}
	cookieSession, _ := cookieStore.New(nil, CookieSessionId)

	ctx.Set(CtxSessionStore, sessionStore)
	ctx.Set(CtxCookie, cookieStore)

	_, err := getSessionFromCookie(ctx)
	if err == nil {
		t.Error("expected error")
	}

	cookieSession.Values[CookieValueSessionId] = "first_session_id"
	_ = cookieStore.Save(nil, nil, cookieSession)

	_, err = getSessionFromCookie(ctx)
	if err == nil {
		t.Error("expected error")
	}

	newSess, _ := sessionStore.NewOrExisting("first_session_id", "")

	sess, err := getSessionFromCookie(ctx)
	if err != nil {
		t.Error("unexpected error", err)
	}

	if sess != newSess {
		t.Error("expected matching sessions", sess, newSess)
	}
}

func TestIdentityMiddleware(t *testing.T) {
	ctx := &_fakeContext{values: make(map[string]interface{})}
	cfg := session.Config{}
	sessionStore := session.NewSessionStore(cfg)
	cookieStore := &_fakeCookieStore{sessions: make(map[string]*sessions.Session)}
	cookieSession, _ := cookieStore.New(nil, CookieSessionId)
	cookieSession.Values[CookieValueSessionId] = "first_session_id"
	_ = cookieStore.Save(nil, nil, cookieSession)

	newSess, _ := sessionStore.NewOrExisting("first_session_id", "")

	ctx.Set(CtxSessionStore, sessionStore)
	ctx.Set(CtxCookie, cookieStore)

	IdentityMiddleware(ctx)(nil)

	sess := ctx.Get(CtxSession)
	if sess != newSess {
		t.Error("expected matching sessions")
	}

	//TODO: user in ctx
}

func _assertPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("The code did not panic")
	}
}
