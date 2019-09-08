package controllers

import (
	"fmt"
	"goairmon/business/data/context"
	"goairmon/business/services/identity"
	"goairmon/site/helper"
	"goairmon/site/models"
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
		loginVM := models.UnmarshalLoginVm(c)

		if err := loginUser(c, identity, loginVM); err != nil {
			view := loadView("auth/login.gohtml", c)
			vm := models.NewContextVm(c, loginVM)
			vm.Errors["general"] = "Failed to log in"

			return view.Execute(c.Response().Writer, vm)
		}

		err := getFlashService(c).PushSuccess(c, "Successfully logged in!")
		if err != nil {
			log.Println(err)
		}

		return c.Redirect(http.StatusSeeOther, "/")
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

func loginUser(c echo.Context, identity *identity.IdentityService, loginVM *models.LoginVm) error {
	dbContext := c.Get(helper.CtxDbContext).(context.DbContext)
	user, err := dbContext.FindUserByName(loginVM.Username)
	if err != nil || !user.CheckPassword(loginVM.Password) {
		return fmt.Errorf("invalid username or password")
	}

	session, err := identity.StartNewSession(c)
	if err != nil {
		return fmt.Errorf("oops! something went wrong")
	}

	session.Values["user_name"] = user.Username

	return nil
}
