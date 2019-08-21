package flash

import (
	"goairmon/site/models"
	"goairmon/site/testhelpers"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
)

func TestFlashStoreInCookie(t *testing.T) {
	cookies := sessions.NewCookieStore([]byte("cookie-encryption-key"))
	ctx := &testhelpers.FakeContext{FakeWriter: httptest.NewRecorder(), Values: make(map[string]interface{})}
	cookieSession, _ := cookies.New(ctx.Request(), "gowebapp_session")
	ctx.Set(CtxCookieSession, cookieSession)
	service := FlashService{}

	if err := service.PushInfo(ctx, "test-info"); err != nil {
		t.Error("unexpected error", err)
	}

	if err := service.PushError(ctx, "test-error"); err != nil {
		t.Error("unnexpected error", err)
	}

	if err := service.PushSuccess(ctx, "test-success"); err != nil {
		t.Error("unnexpected error", err)
	}

	if err := service.PopToContext()(testhelpers.EmptyHandler)(ctx); err != nil {
		t.Error("unexpected error", err)
	}

	flashBag := ctx.Get(CtxFlashMessages).(*models.FlashBag)

	if len(flashBag.Info) != 1 || flashBag.Info[0] != "test-info" {
		t.Error("unexpected flash bag values")
	}

	if len(flashBag.Success) != 1 || flashBag.Success[0] != "test-success" {
		t.Error("unexpected flash bag values")
	}

	if len(flashBag.Error) != 1 || flashBag.Error[0] != "test-error" {
		t.Error("unexpected flash bag values")
	}
}
