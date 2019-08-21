package flash

import (
	"goairmon/site/helper"
	"goairmon/site/models"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
)

const (
	CookieFlashBagKey = "flash_bag"
	CtxCookieSession  = helper.CtxCookieSession
	CtxFlashMessages  = helper.CtxFlashMessages
)

type FlashService struct {
	//
}

func (f *FlashService) PopToContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookieSession := f.mustGetCookieSession(c)
			flashBag := f.getBagFromCookie(cookieSession, true)
			_ = cookieSession.Save(c.Request(), c.Response().Writer)
			c.Set(CtxFlashMessages, flashBag)

			return next(c)
		}
	}
}

func (f *FlashService) PushSuccess(c echo.Context, msg string) error {
	session := f.mustGetCookieSession(c)
	flashBag := f.getBagFromCookie(session, false)
	flashBag.Success = append(flashBag.Success, msg)

	return f.saveBagToCookie(flashBag, session, c)
}

func (f *FlashService) PushInfo(c echo.Context, msg string) error {
	session := f.mustGetCookieSession(c)
	flashBag := f.getBagFromCookie(session, false)
	flashBag.Info = append(flashBag.Info, msg)

	return f.saveBagToCookie(flashBag, session, c)
}

func (f *FlashService) PushError(c echo.Context, msg string) error {
	session := f.mustGetCookieSession(c)
	flashBag := f.getBagFromCookie(session, false)
	flashBag.Error = append(flashBag.Error, msg)

	return f.saveBagToCookie(flashBag, session, c)
}

func (f *FlashService) getBagFromCookie(cookieSession *sessions.Session, delete bool) *models.FlashBag {
	raw := f.getCookieValue(cookieSession, CookieFlashBagKey, delete)

	bag := &models.FlashBag{}
	_ = bag.Decode(raw)

	return bag
}

func (f *FlashService) saveBagToCookie(bag *models.FlashBag, cookieSession *sessions.Session, c echo.Context) error {
	cookieSession.Values[CookieFlashBagKey] = bag.Encode()

	return cookieSession.Save(c.Request(), c.Response().Writer)
}

func (f *FlashService) mustGetCookieSession(c echo.Context) *sessions.Session {
	return c.Get(CtxCookieSession).(*sessions.Session)
}

func (f *FlashService) getCookieValue(cookieSession *sessions.Session, key string, delete bool) string {
	var value string
	if encoded, ok := cookieSession.Values[key].(string); ok {
		value = encoded
	} else {
		value = ""
	}

	if delete {
		cookieSession.Values[key] = nil
	}

	return value
}
