package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thitiphum-bluesage/assessment-tax/applications/services"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/schemas"
	"github.com/thitiphum-bluesage/assessment-tax/utilities"
)

type TaxController struct {
	taxService services.TaxServiceInterface
}

func NewTaxController(service services.TaxServiceInterface) *TaxController {
	return &TaxController{
		taxService: service,
	}
}

func (tc *TaxController) CalculateTax(c echo.Context) error {
	var req schemas.TaxCalculationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input data")
	}

	if err := utilities.ValidateTaxCalculationRequest(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tax, err := tc.taxService.CalculateTax(*req.TotalIncome, *req.WHT, req.Allowances)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := schemas.TaxCalculationResponse{
		Tax: tax,
	}
	return c.JSON(http.StatusOK, response)
}
