package controllers

import (
	"html/template"
	"os"
	"strings"

	"github.com/labstack/echo"
)

func NewAuthController() Controller {

}

type AuthController struct {
}

func (a *AuthController) Group(server echo.Echo) *echo.Group {
	group := server.Group("auth")
	group.GET("/login", func(c echo.Context) error {
		view := loadView("auth/login.tmpl", c)

	})

	return group
}

func siteRoot() string {
	for i := 0; i < 10; i++ {
		dir, _ := os.Getwd()
		if strings.HasSuffix(dir, "/site") {
			return dir
		}

		dir = "../" + dir
	}

	panic("failed to find site root")
}

func loadView(viewPath string, c echo.Context) *template.Template {
	template.ParseFiles(viewPath)
}
