package models

import (
	"goairmon/business/services/session"
	"goairmon/site/helper"

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
	userName := ""
	if sess != nil {
		userName = sess.Values["user_name"]
	}

	return &ContextVm{
		Session:   sess,
		UserName:  userName,
		ViewModel: viewModel,
		Errors:    ErrorBag{},
		FlashBag:  flashBag,
		Csrf:      csrfToken,
	}
}

type ContextVm struct {
	Session   *session.Session
	UserName  string
	ViewModel interface{}
	Errors    ErrorBag
	FlashBag  *FlashBag
	Csrf      string
}
