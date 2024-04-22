package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thitiphum-bluesage/assessment-tax/config"
)

func BasicAuth(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username, password, ok := c.Request().BasicAuth()
			if !ok || username != cfg.AdminUser || password != cfg.AdminPass {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: Incorrect credentials")
			}
			return next(c)
		}
	}
}
