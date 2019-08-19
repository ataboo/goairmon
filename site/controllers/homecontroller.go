package controllers

import (
	"goairmon/site/models"
	"goairmon/site/services/identity"

	"github.com/labstack/echo"
)

func HomeController(server *echo.Echo, identity *identity.IdentityService) *echo.Group {
	group := server.Group("", identity.RedirectUsersWithoutSession("/auth/login"))
	group.GET("/", func(c echo.Context) error {
		view := loadView("home/index.gohtml", c)

		return view.Execute(c.Response().Writer, models.NewContextVm(c, nil))
	})

	return group
}
