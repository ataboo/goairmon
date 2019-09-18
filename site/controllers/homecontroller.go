package controllers

import (
	"goairmon/business/services/identity"
	"goairmon/site/models"

	"github.com/labstack/echo"
)

func HomeController(server *echo.Echo, identity *identity.IdentityService) *echo.Group {
	group := server.Group("")
	group.GET("/", func(c echo.Context) error {
		view := loadView("home/index.gohtml", c)

		return view.Execute(c.Response().Writer, models.NewContextVm(c, nil))
	}, identity.RedirectUsersWithoutSession("/auth/login"))

	return group
}
