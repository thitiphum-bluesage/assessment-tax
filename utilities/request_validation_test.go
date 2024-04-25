package utilities

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/schemas"
)

func TestValidateUpdatePersonalDeductionRequest(t *testing.T) {
	validAmount := 50000.0
	tooLowAmount := 5000.0
	tooHighAmount := 150000.0

	tests := []struct {
		name     string
		input    *schemas.UpdatePersonalDeductionRequest
		expected string
	}{
		{"Valid Amount", &schemas.UpdatePersonalDeductionRequest{Amount: &validAmount}, ""},
		{"Amount is nil", &schemas.UpdatePersonalDeductionRequest{Amount: nil}, "amount is required"},
		{"Amount too low", &schemas.UpdatePersonalDeductionRequest{Amount: &tooLowAmount}, "amount must be between 10,000 and 100,000"},
		{"Amount too high", &schemas.UpdatePersonalDeductionRequest{Amount: &tooHighAmount}, "amount must be between 10,000 and 100,000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdatePersonalDeductionRequest(tt.input)
			if err != nil {
				assert.Equal(t, tt.expected, err.Error(), "Expected error message to match")
			} else {
				assert.Empty(t, tt.expected, "Expected no error")
			}
		})
	}
}

func TestValidateUpdateKReceiptRequest(t *testing.T) {
	validAmount := 50000.0
	tooLowAmount := 0.0
	tooHighAmount := 150000.0

	tests := []struct {
		name     string
		input    *schemas.UpdateKReceiptRequest
		expected string
	}{
		{"Valid Amount", &schemas.UpdateKReceiptRequest{Amount: &validAmount}, ""},
		{"Amount is nil", &schemas.UpdateKReceiptRequest{Amount: nil}, "amount for k-receipt is required"},
		{"Amount too low", &schemas.UpdateKReceiptRequest{Amount: &tooLowAmount}, "amount for k-receipt must be between 1 and 100,000"},
		{"Amount too high", &schemas.UpdateKReceiptRequest{Amount: &tooHighAmount}, "amount for k-receipt must be between 1 and 100,000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdateKReceiptRequest(tt.input)
			if err != nil {
				assert.Equal(t, tt.expected, err.Error(), "Expected error message to match")
			} else {
				assert.Empty(t, tt.expected, "Expected no error")
			}
		})
	}
}


func TestValidateTaxCalculationRequest(t *testing.T) {
    positiveIncome := 50000.0
    positiveWHT := 3000.0
    negativeWHT := -100.0
    higherWHT := 60000.0

    tests := []struct {
        name     string
        input    *schemas.TaxCalculationRequest
        expected string
    }{
        {"Valid Input with Donation", &schemas.TaxCalculationRequest{
            TotalIncome: &positiveIncome,
            WHT:         &positiveWHT,
            Allowances:  []schemas.Allowance{{AllowanceType: "donation", Amount: 500}},
        }, ""},
        {"Valid Input with K-Receipt", &schemas.TaxCalculationRequest{
            TotalIncome: &positiveIncome,
            WHT:         &positiveWHT,
            Allowances:  []schemas.Allowance{{AllowanceType: "k-receipt", Amount: 1000}},
        }, ""},
        {"Invalid Allowance Type", &schemas.TaxCalculationRequest{
            TotalIncome: &positiveIncome,
            WHT:         &positiveWHT,
            Allowances:  []schemas.Allowance{{AllowanceType: "invalid_type", Amount: 500}},
        }, "validation errors: Invalid allowance type: invalid_type. Allowed types are 'donation' and 'k-receipt'."},
        {"Negative Allowance Amount", &schemas.TaxCalculationRequest{
            TotalIncome: &positiveIncome,
            WHT:         &positiveWHT,
            Allowances:  []schemas.Allowance{{AllowanceType: "donation", Amount: -500}},
        }, "validation errors: Amount for donation must be non-negative"},
        {"Duplicate Allowances", &schemas.TaxCalculationRequest{
            TotalIncome: &positiveIncome,
            WHT:         &positiveWHT,
            Allowances: []schemas.Allowance{
                {AllowanceType: "donation", Amount: 500},
                {AllowanceType: "donation", Amount: 300},
            },
        }, "validation errors: Only one donation allowance can be included"},
        {"TotalIncome is nil", &schemas.TaxCalculationRequest{
            WHT: &positiveWHT,
        }, "validation errors: TotalIncome is required"},
        {"WHT is nil", &schemas.TaxCalculationRequest{
            TotalIncome: &positiveIncome,
        }, "validation errors: WHT is required"},
        {"WHT is negative", &schemas.TaxCalculationRequest{
            TotalIncome: &positiveIncome,
            WHT:         &negativeWHT,
        }, "validation errors: WHT must be non-negative"},
        {"WHT greater than TotalIncome", &schemas.TaxCalculationRequest{
            TotalIncome: &positiveIncome,
            WHT:         &higherWHT,
        }, "validation errors: WHT cannot be greater than TotalIncome"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateTaxCalculationRequest(tt.input)
            if err != nil {
                assert.Equal(t, tt.expected, err.Error(), "Expected error message to match")
            } else {
                assert.Empty(t, tt.expected, "Expected no error")
            }
        })
    }
}




func TestValidateCSVTaxRecords(t *testing.T) {
	tests := []struct {
		name     string
		input    []schemas.CSVObjectFormat
		expected string
	}{
		{"Valid Records", []schemas.CSVObjectFormat{
			{TotalIncome: 50000, WHT: 3000, Donation: 500, KReceipt: 200},
			{TotalIncome: 75000, WHT: 5000, Donation: 1000, KReceipt: 500},
		}, ""},
		{"Negative TotalIncome", []schemas.CSVObjectFormat{
			{TotalIncome: -100, WHT: 3000, Donation: 500, KReceipt: 200},
		}, "validation errors: Record 1: TotalIncome must be non-negative; Record 1: WHT cannot be greater than TotalIncome"},
		{"Negative WHT", []schemas.CSVObjectFormat{
			{TotalIncome: 50000, WHT: -100, Donation: 500, KReceipt: 200},
		}, "validation errors: Record 1: WHT must be non-negative"},
		{"Negative Donation", []schemas.CSVObjectFormat{
			{TotalIncome: 50000, WHT: 3000, Donation: -500, KReceipt: 200},
		}, "validation errors: Record 1: Donation must be non-negative"},
		{"Negative KReceipt", []schemas.CSVObjectFormat{
			{TotalIncome: 50000, WHT: 3000, Donation: 500, KReceipt: -200},
		}, "validation errors: Record 1: KReceipt must be non-negative"},
		{"WHT greater than TotalIncome", []schemas.CSVObjectFormat{
			{TotalIncome: 5000, WHT: 6000, Donation: 500, KReceipt: 200},
		}, "validation errors: Record 1: WHT cannot be greater than TotalIncome"},
		{"Invalid TotalIncome Type", []schemas.CSVObjectFormat{
			{TotalIncome: math.NaN(), WHT: 3000, Donation: 500, KReceipt: 200},
		}, "validation errors: Record 1: TotalIncome must be a valid float"},
		{"Invalid Donation Type", []schemas.CSVObjectFormat{
			{TotalIncome: 50000, WHT: 3000, Donation: math.NaN(), KReceipt: 200},
		}, "validation errors: Record 1: Donation must be a valid float"},
		{"Multiple Validation Errors", []schemas.CSVObjectFormat{
			{TotalIncome: -100, WHT: -200, Donation: -300, KReceipt: -400},
		}, "validation errors: Record 1: TotalIncome must be non-negative; Record 1: WHT must be non-negative; Record 1: Donation must be non-negative; Record 1: KReceipt must be non-negative"},
		{"Multiple Records with Errors", []schemas.CSVObjectFormat{
			{TotalIncome: 50000, WHT: 3000, Donation: 500, KReceipt: 200},
			{TotalIncome: -100, WHT: 3000, Donation: 500, KReceipt: 200},
			{TotalIncome: 60000, WHT: -200, Donation: 1000, KReceipt: 300},
		}, "validation errors: Record 2: TotalIncome must be non-negative; Record 2: WHT cannot be greater than TotalIncome; Record 3: WHT must be non-negative"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCSVTaxRecords(tt.input)
			if err != nil {
				assert.Equal(t, tt.expected, err.Error(), "Expected error message to match")
			} else {
				assert.Empty(t, tt.expected, "Expected no error")
			}
		})
	}
}
