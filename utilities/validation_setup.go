package utilities

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator struct wraps a validator.Validate to use with Echo framework.
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator creates and returns a new CustomValidator.
func NewValidator() echo.Validator {
	v := validator.New()
	return &CustomValidator{validator: v}
}

// Validate executes the validator on the given interface.
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
