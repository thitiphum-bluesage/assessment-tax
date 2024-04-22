package schemas

import (
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

type UpdatePersonalDeductionRequest struct {
	Amount float64 `json:"amount" validate:"required,gte=10000,lte=100000"`
}

type UpdatePersonalDeductionResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type UpdateKReceiptRequest struct {
	Amount float64 `json:"amount" validate:"required,gte=1,lte=100000"`
}

type UpdateKReceiptResponse struct {
	KReceipt float64 `json:"kReceipt"`
}
