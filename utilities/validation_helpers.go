package utilities

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) string {
	var errMessages []string
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			errMessages = append(errMessages, formatFieldError(e))
		}
		return strings.Join(errMessages, ", ")
	}
	return err.Error()
}

func formatFieldError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "max":
		return fmt.Sprintf("%s cannot be more than %s", e.Field(), e.Param())
	case "min":
		return fmt.Sprintf("%s must be at least %s", e.Field(), e.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", e.Field())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", e.Field(), e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s is not valid", e.Field())
	}
}
