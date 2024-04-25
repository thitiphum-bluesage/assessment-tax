package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thitiphum-bluesage/assessment-tax/config"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/endpoints/controllers"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/middleware"
)

func NewRouter(e *echo.Echo, taxControllerr *controllers.TaxController, adminController *controllers.AdminController, cfg *config.Config) {

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	// Group for tax-related routes
	taxGroup := e.Group("/tax")
	taxGroup.POST("/calculations", taxControllerr.CalculateDetailedTax)
	taxGroup.POST("/calculations/upload-csv", taxControllerr.CalculateCSVTax)

	// Group for admin-related routes
	adminGroup := e.Group("/admin")
	adminGroup.Use(middleware.BasicAuth(cfg))
	adminGroup.POST("/deductions/personal", adminController.UpdatePersonalDeduction)
	adminGroup.POST("/deductions/k-receipt", adminController.UpdateKReceiptDeduction)
}
