package tax

import (
	"math"

	"github.com/thitiphum-bluesage/assessment-tax/interfaces/schemas"
	"github.com/thitiphum-bluesage/assessment-tax/utilities"
)

func getTaxBrackets() []struct {
	UpperBound float64
	TaxRate    float64
} {
	return []struct {
		UpperBound float64
		TaxRate    float64
	}{
		{150000, 0},
		{500000, 0.1},
		{1000000, 0.15},
		{2000000, 0.2},
		{math.MaxFloat64, 0.35},
	}
}

func calculateProgressiveTax(income float64) float64 {
	brackets := getTaxBrackets()
	tax := 0.0
	previousUpperBound := 0.0

	for _, bracket := range brackets {
		if income > bracket.UpperBound {
			taxAmount := (bracket.UpperBound - previousUpperBound) * bracket.TaxRate
			taxAmount = utilities.FormatToTwoDecimals(taxAmount)
			tax += taxAmount
			previousUpperBound = bracket.UpperBound
			continue
		}
		taxAmount := (income - previousUpperBound) * bracket.TaxRate
		taxAmount = utilities.FormatToTwoDecimals(taxAmount)
		tax += taxAmount
		break
	}
	return tax
}

func calculateProgressiveTaxWithDetails(income float64) ([]schemas.TaxLevel, float64) {
	brackets := getTaxBrackets()

	detailResponse := []schemas.TaxLevel{
		{Level: "0-150,000", Tax: 0.0},
		{Level: "150,001-500,000", Tax: 0.0},
		{Level: "500,001-1,000,000", Tax: 0.0},
		{Level: "1,000,001-2,000,000", Tax: 0.0},
		{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
	}

	tax := 0.0
	previousUpperBound := 0.0

	for i, bracket := range brackets {
		if income > bracket.UpperBound {
			taxAmount := (bracket.UpperBound - previousUpperBound) * bracket.TaxRate
			taxAmount = utilities.FormatToTwoDecimals(taxAmount)
			tax += taxAmount
			previousUpperBound = bracket.UpperBound
			detailResponse[i].Tax = taxAmount
			continue
		}
		taxAmount := (income - previousUpperBound) * bracket.TaxRate
		taxAmount = utilities.FormatToTwoDecimals(taxAmount)
		tax += taxAmount
		detailResponse[i].Tax = taxAmount
		break
	}

	return detailResponse, tax
}
