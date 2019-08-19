package controllers

import (
	"fmt"
	"goairmon/site/helper"
	"goairmon/site/services/flash"
	"goairmon/site/services/viewloader"
	"html/template"

	"github.com/labstack/echo"
)

const (
	CtxFlashServiceKey = helper.CtxFlashServiceKey
)

func loadView(viewPath string, c echo.Context) *template.Template {
	provider, ok := c.Get(viewloader.CtxKey).(*viewloader.ViewLoader)
	if !ok || provider == nil {
		panic(fmt.Sprintf("failed to load view: %s", viewPath))
	}

	return provider.LoadView(viewPath, c)
}

func getFlashService(c echo.Context) *flash.FlashService {
	return c.Get(CtxFlashServiceKey).(*flash.FlashService)
}
