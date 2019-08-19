package controllers

import (
	"goairmon/site/models"
	"goairmon/site/services/identity"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func AuthController(server *echo.Echo, identity *identity.IdentityService) *echo.Group {
	group := server.Group("auth")
	group.GET("/login", func(c echo.Context) error {
		view := loadView("auth/login.gohtml", c)

		return view.Execute(c.Response().Writer, models.NewContextVm(c, nil))
	}, identity.RedirectUsersWithSession("/"))

	group.POST("/login", func(c echo.Context) error {
		loginVm := models.UnmarshalLoginVm(c)

		if loginVm.Username == "ataboo" && loginVm.Password == "asdfasdf" {
			_ = identity.StartNewSession(c)
			err := getFlashService(c).PushSuccess(c, "Successfully logged in!")
			if err != nil {
				log.Println(err)
			}

			return c.Redirect(http.StatusSeeOther, "/")
		}

		view := loadView("auth/login.gohtml", c)
		vm := models.NewContextVm(c, loginVm)
		vm.Errors["general"] = "Invalid username or password"

		return view.Execute(c.Response().Writer, vm)
	}, identity.RedirectUsersWithSession("/"))

	group.POST("/logout", func(c echo.Context) error {
		_ = identity.EndSession(c)

		returnPath := c.FormValue("referer")
		if returnPath == "" {
			returnPath = "/"
		}

		return c.Redirect(http.StatusSeeOther, returnPath)
	}, identity.RedirectUsersWithoutSession("/"))

	return group
}
