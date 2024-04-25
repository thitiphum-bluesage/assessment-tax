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

	netTax, taxRefund, err := tc.taxService.CalculateTax(*req.TotalIncome, *req.WHT, req.Allowances)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if taxRefund > 0 {
		response := schemas.TaxCalculationRefundResponse{
			TaxRefund: taxRefund,
		}
		return c.JSON(http.StatusOK, response)
	}

	response := schemas.TaxCalculationResponse{
		Tax: netTax,
	}
	return c.JSON(http.StatusOK, response)
}

func (tc *TaxController) CalculateDetailedTax(c echo.Context) error {
	var req schemas.TaxCalculationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input data")
	}

	if err := utilities.ValidateTaxCalculationRequest(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	taxLevel, netTax, taxRefund, err := tc.taxService.CalculateDetailedTax(*req.TotalIncome, *req.WHT, req.Allowances)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if taxRefund > 0 {
		response := schemas.TaxCalculationRefundResponse{
			TaxRefund: taxRefund,
		}
		return c.JSON(http.StatusOK, response)
	}

	response := schemas.DetailedTaxCalculationResponse{
		Tax:      netTax,
		TaxLevel: taxLevel,
	}
	return c.JSON(http.StatusOK, response)
}
