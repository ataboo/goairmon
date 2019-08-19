package models

import (
	"goairmon/site/helper"
	"goairmon/site/services/session"

	"github.com/labstack/echo"
)

const (
	CtxKeySession    = helper.CtxKeySession
	CtxFlashMessages = helper.CtxFlashMessages
)

func NewContextVm(c echo.Context, viewModel interface{}) *ContextVm {
	// If this gets carried away, make me a factory service
	sess, _ := c.Get(CtxKeySession).(*session.Session)
	flashBag, _ := c.Get(CtxFlashMessages).(*FlashBag)
	csrfToken := c.Get("csrf").(string)

	return &ContextVm{
		Session:   sess,
		ViewModel: viewModel,
		Errors:    ErrorBag{},
		FlashBag:  flashBag,
		Csrf:      csrfToken,
	}
}

type ContextVm struct {
	Session   *session.Session
	ViewModel interface{}
	Errors    ErrorBag
	FlashBag  *FlashBag
	Csrf      string
}
