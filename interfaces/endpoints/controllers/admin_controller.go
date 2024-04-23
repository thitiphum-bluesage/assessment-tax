package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thitiphum-bluesage/assessment-tax/applications/services"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/schemas"
	"github.com/thitiphum-bluesage/assessment-tax/utilities"
)

type AdminController struct {
	service services.AdminServiceInterface
}

func NewAdminController(service services.AdminServiceInterface) *AdminController {
	return &AdminController{
		service: service,
	}
}

func (ac *AdminController) UpdatePersonalDeduction(c echo.Context) error {
	var req schemas.UpdatePersonalDeductionRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input data")
	}

	if err := utilities.ValidateUpdatePersonalDeductionRequest(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ac.service.UpdatePersonalDeduction(*req.Amount); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, schemas.UpdatePersonalDeductionResponse{
		PersonalDeduction: *req.Amount,
	})
}

func (ac *AdminController) UpdateKReceiptDeduction(c echo.Context) error {
	var req schemas.UpdateKReceiptRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input data")
	}

	if err := utilities.ValidateUpdateKReceiptRequest(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ac.service.UpdateKReceiptDeductionMax(*req.Amount); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, schemas.UpdateKReceiptResponse{KReceipt: *req.Amount})
}
