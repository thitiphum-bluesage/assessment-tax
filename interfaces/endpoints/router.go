package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/endpoints/controllers"
)

func NewRouter(e *echo.Echo, adminController *controllers.AdminController) {

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})


    e.POST("/admin/deductions/personal", adminController.UpdatePersonalDeduction)
}
