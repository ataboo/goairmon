package controllers

import "github.com/labstack/echo"

type Controller interface {
	Group() *echo.Group
}
