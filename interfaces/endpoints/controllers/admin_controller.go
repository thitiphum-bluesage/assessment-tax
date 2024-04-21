package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thitiphum-bluesage/assessment-tax/applications/services"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/schemas"
)

type AdminController struct {
    service *services.AdminService
}

func NewAdminController(service *services.AdminService) *AdminController {
    return &AdminController{
        service: service,
    }
}

func (ac *AdminController) UpdatePersonalDeduction(c echo.Context) error {
    var req schemas.UpdatePersonalDeductionRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid input data")
    }

    if err := schemas.Validate.Struct(req); err != nil {
        // Return validation errors
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    if err := ac.service.UpdatePersonalDeduction(req.Amount); err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

    return c.JSON(http.StatusOK, schemas.UpdatePersonalDeductionResponse{
        PersonalDeduction: req.Amount,
    })
}