package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thitiphum-bluesage/assessment-tax/config"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/endpoints/controllers"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/middleware"
)

func NewRouter(e *echo.Echo, adminController *controllers.AdminController, cfg *config.Config) {

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	adminGroup := e.Group("/admin")
	adminGroup.Use(middleware.BasicAuth(cfg))

	adminGroup.POST("/deductions/personal", adminController.UpdatePersonalDeduction)
}
