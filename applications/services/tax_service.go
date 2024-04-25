package services

import (
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/schemas"
)

type TaxServiceInterface interface {
	CalculateTax(income float64, wht float64, allowances []schemas.Allowance) (float64, float64, error)
	CalculateDetailedTax(income float64, wht float64, allowances []schemas.Allowance) ([]schemas.TaxLevel, float64, float64, error)
	CalculateTaxFromCSV(records []schemas.CSVObjectFormat) (schemas.CSVResponse, error)
}
