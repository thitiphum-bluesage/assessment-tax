package services

import (
	"math"

	"github.com/thitiphum-bluesage/assessment-tax/infrastructure/repository"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/schemas"
	"github.com/thitiphum-bluesage/assessment-tax/utilities"
)

type taxService struct {
	taxRepo repository.TaxDeductionConfigRepositoryInterface
}

func NewTaxService(taxRepo repository.TaxDeductionConfigRepositoryInterface) TaxServiceInterface {
	return &taxService{
		taxRepo: taxRepo,
	}
}

func (s *taxService) CalculateTax(income float64, wht float64, allowances []schemas.Allowance) (float64, float64, error) {
	config, err := s.taxRepo.GetConfig()
	if err != nil {
		return 0, 0, err
	}

	allowancesDeduction := config.PersonalDeduction
	for _, allowance := range allowances {
		if allowance.AllowanceType == "donation" {
			if allowance.Amount > config.DonationDeductionMax {
				allowancesDeduction += config.DonationDeductionMax
			} else {
				allowancesDeduction += allowance.Amount
			}
		} else {
			if allowance.Amount > config.KReceiptDeductionMax {
				allowancesDeduction += config.KReceiptDeductionMax
			} else {
				allowancesDeduction += allowance.Amount
			}
		}
	}

	incomeAfterDeduct := income - allowancesDeduction
	if incomeAfterDeduct < 0 {
		incomeAfterDeduct = 0
	}

	tax := calculateProgressiveTax(incomeAfterDeduct)

	netTax := tax - wht
	taxRefund := 0.0
	if netTax < 0 {
		taxRefund = -netTax // Calculate refund as the negative of a negative tax value
		netTax = 0
	}

	return netTax, taxRefund, nil
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

func (s *taxService) CalculateDetailedTax(income, wht float64, allowances []schemas.Allowance) ([]schemas.TaxLevel, float64, float64, error) {
	config, err := s.taxRepo.GetConfig()
	if err != nil {
		return nil, 0, 0, err
	}

	allowancesDeduction := config.PersonalDeduction
	for _, allowance := range allowances {
		if allowance.AllowanceType == "donation" {
			if allowance.Amount > config.DonationDeductionMax {
				allowancesDeduction += config.DonationDeductionMax
			} else {
				allowancesDeduction += allowance.Amount
			}
		} else {
			if allowance.Amount > config.KReceiptDeductionMax {
				allowancesDeduction += config.KReceiptDeductionMax
			} else {
				allowancesDeduction += allowance.Amount
			}
		}
	}

	incomeAfterDeduct := income - allowancesDeduction
	if incomeAfterDeduct < 0 {
		incomeAfterDeduct = 0
	}

	taxLevels, tax := calculateProgressiveTaxWithDetails(incomeAfterDeduct)

	netTax := tax - wht
	taxRefund := 0.0
	if netTax < 0 {
		taxRefund = -netTax // Calculate refund as the negative of a negative tax value
		netTax = 0
	}

	return taxLevels, netTax, taxRefund, nil

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

func (s *taxService) CalculateTaxFromCSV(records []schemas.CSVObjectFormat) (schemas.CSVResponse, error) {
	config, err := s.taxRepo.GetConfig()
	if err != nil {
		return schemas.CSVResponse{}, err
	}

	var response schemas.CSVResponse

	for _, record := range records {
		totalIncome := record.TotalIncome
		wht := record.WHT
		donation := record.Donation
		k_receipt := record.KReceipt

		if donation > config.DonationDeductionMax {
			donation = config.DonationDeductionMax
		}
		if k_receipt > config.KReceiptDeductionMax {
			k_receipt = config.KReceiptDeductionMax
		}

		totalIncomeAfterDeduct := totalIncome - (donation + k_receipt + config.PersonalDeduction)

		if totalIncomeAfterDeduct < 0 {
			totalIncomeAfterDeduct = 0
		}

		tax := calculateProgressiveTax(totalIncomeAfterDeduct)
		netTax := tax - wht
		if netTax < 0 {
			response.Taxes = append(response.Taxes, schemas.CSVResponseMember{
				TotalIncome: totalIncome,
				TaxRefund:   -netTax,
			})
		} else {
			response.Taxes = append(response.Taxes, schemas.CSVResponseMember{
				TotalIncome: totalIncome,
				Tax:         netTax,
			})
		}
	}
	return response, nil
}

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
