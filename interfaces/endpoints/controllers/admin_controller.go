package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thitiphum-bluesage/assessment-tax/applications/services/admin"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/schemas"
	"github.com/thitiphum-bluesage/assessment-tax/utilities"
)

type AdminController struct {
	service admin.AdminServiceInterface
}

func NewAdminController(service admin.AdminServiceInterface) *AdminController {
	return &AdminController{
		service: service,
	}
}

// UpdatePersonalDeduction updates the personal tax deduction amount
// @Summary Update personal deduction
// @Description Update the personal deduction for a tax payer
// @Tags admin
// @Accept json
// @Produce json
// @Param request body schemas.UpdatePersonalDeductionRequest true "Update Personal Deduction Request"
// @Success 200 {object} schemas.UpdatePersonalDeductionResponse
// @Failure 400 {object} schemas.ErrorResponse "Invalid input data"
// @Failure 500 {object} schemas.ErrorResponse "Internal Server Error"
// @Security basicAuth
// @Router /admin/deductions/personal [post]
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

// UpdateKReceiptDeduction updates the K receipt deduction amount
// @Summary Update K receipt deduction
// @Description Update the K receipt deduction for a tax payer
// @Tags admin
// @Accept json
// @Produce json
// @Param request body schemas.UpdateKReceiptRequest true "Update K Receipt Deduction Request"
// @Success 200 {object} schemas.UpdateKReceiptResponse
// @Failure 400 {object} schemas.ErrorResponse "Invalid input data"
// @Failure 500 {object} schemas.ErrorResponse "Internal Server Error"
// @Security basicAuth
// @Router /admin/deductions/k-receipt [post]
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
