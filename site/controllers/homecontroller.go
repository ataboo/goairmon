package controllers

import (
	"fmt"
	"goairmon/business/data/context"
	"goairmon/business/services/identity"
	"goairmon/site/helper"
	"goairmon/site/models"

	"github.com/labstack/echo"
)

func HomeController(server *echo.Echo, identity *identity.IdentityService) *echo.Group {
	group := server.Group("")
	group.GET("/", func(c echo.Context) error {
		view := loadView("home/index.gohtml", c)

		dbContext := c.Get(helper.CtxDbContext).(context.DbContext)
		points, err := dbContext.GetSensorPoints(24 * 60)
		if err != nil {
			return fmt.Errorf("failed to get sensor points", err)
		}

		return view.Execute(c.Response().Writer, models.NewContextVm(c, models.GraphVm{SensorPoints: points}))
	}, identity.RedirectUsersWithoutSession("/auth/login"))

	return group
}
