package flash

import (
	"goairmon/site/context"
	"goairmon/site/models"
	"goairmon/site/services/identity"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

const (
	CookieFlashSuccess = "flash_success"
	CookieFlashInfo    = "flash_info"
	CookieFlashDanger  = "flash_danger"
)

type FlashService struct {
	//
}

func (f *FlashService) PopToContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			flashBag := f.getBagFromCookie(f.mustGetCookieSession(c), true)
			c.Set(context.CtxFlashKey, flashBag)

			return next(c)
		}
	}
}

func (f *FlashService) PushSuccess(c echo.Context, msg string) error {
	session := f.mustGetCookieSession(c)
	flashBag := f.getBagFromCookie(session, false)
	flashBag.Success = append(flashBag.Success, msg)

	return session.Save(c.Request(), c.Response().Writer)
}

func (f *FlashService) PushInfo(c echo.Context, msg string) error {
	session := f.mustGetCookieSession(c)
	flashBag := f.getBagFromCookie(session, false)
	flashBag.Info = append(flashBag.Info, msg)

	return session.Save(c.Request(), c.Response().Writer)
}

func (f *FlashService) PushError(c echo.Context, msg string) error {
	session := f.mustGetCookieSession(c)
	flashBag := f.getBagFromCookie(session, false)
	flashBag.Error = append(flashBag.Error, msg)

	return session.Save(c.Request(), c.Response().Writer)
}

func (f *FlashService) getBagFromCookie(cookieSession *sessions.Session, delete bool) *models.FlashBag {
	bag := &models.FlashBag{
		Success: f.getCookieValue(cookieSession, CookieFlashSuccess, delete),
		Info:    f.getCookieValue(cookieSession, CookieFlashInfo, delete),
		Error:   f.getCookieValue(cookieSession, CookieFlashDanger, delete),
	}

	return bag
}

func (f *FlashService) mustGetCookieSession(c echo.Context) *sessions.Session {
	return c.Get(identity.CtxCookieSession).(*sessions.Session)
}

func (f *FlashService) getCookieValue(cookieSession *sessions.Session, key string, delete bool) []string {
	var value []string
	if msgs, ok := cookieSession.Values[key].([]string); ok {
		value = msgs
	} else {
		value = []string{}
	}

	if delete {
		cookieSession.Values[key] = nil
	}

	return value
}
