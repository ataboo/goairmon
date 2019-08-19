package models

import (
	"goairmon/site/context"
	"goairmon/site/services/identity"
	"goairmon/site/services/session"

	"github.com/labstack/echo"
)

func NewContextVm(c echo.Context, viewModel interface{}) *ContextVm {
	sess, _ := c.Get(identity.CtxKeySession).(*session.Session)
	flashBag, _ := c.Get(context.CtxFlashKey).(*FlashBag)

	return &ContextVm{
		Session:   sess,
		ViewModel: viewModel,
		Errors:    ErrorBag{},
		FlashBag:  flashBag,
	}
}

type ContextVm struct {
	Session   *session.Session
	ViewModel interface{}
	Errors    ErrorBag
	FlashBag  *FlashBag
}
