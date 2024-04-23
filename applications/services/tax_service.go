package services

import (
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/schemas"
)

type TaxServiceInterface interface {
	CalculateTax(income float64, wht float64, allowances []schemas.Allowance) (float64, error)
}
