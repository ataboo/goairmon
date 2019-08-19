package models

import "github.com/labstack/echo"

type LoginVm struct {
	Username string
	Password string
}

func UnmarshalLoginVm(c echo.Context) *LoginVm {
	return &LoginVm{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
	}
}
