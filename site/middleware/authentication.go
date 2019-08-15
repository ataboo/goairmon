package middleware

import "github.com/labstack/echo"

func RequireSession(c echo.Context) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Get("session")

			return nil
		}
	}
}
