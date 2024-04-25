package tax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thitiphum-bluesage/assessment-tax/domains"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/schemas"
)

type MockTaxRepo struct {
	mock.Mock
}

func (m *MockTaxRepo) GetConfig() (*domains.TaxDeductionConfig, error) {
	args := m.Called()
	if config, ok := args.Get(0).(*domains.TaxDeductionConfig); ok {
		return config, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaxRepo) UpdatePersonalDeduction(amount float64) error {
	args := m.Called(amount)
	return args.Error(0)
}

func (m *MockTaxRepo) UpdateKReceiptDeductionMax(amount float64) error {
	args := m.Called(amount)
	return args.Error(0)
}

func TestCalculateProgressiveTax(t *testing.T) {
	tests := []struct {
		name   string
		income float64
		want   float64
	}{
		{"Zero income", 0, 0},
		{"Boundary of first bracket", 150000, 0},
		{"Middle of second bracket", 300000, 15000},
		{"Boundary of second bracket", 500000, 35000},
		{"Above all brackets", 2500000, 485000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateProgressiveTax(tt.income); got != tt.want {
				t.Errorf("calculateProgressiveTax(%f) = %f, want %f", tt.income, got, tt.want)
			}
		})
	}
}

func TestCalculateTaxFromCSV(t *testing.T) {
	mockRepo := new(MockTaxRepo)
	service := NewTaxService(mockRepo)
	config := &domains.TaxDeductionConfig{
		PersonalDeduction:    60000,
		DonationDeductionMax: 100000,
		KReceiptDeductionMax: 50000,
	}

	mockRepo.On("GetConfig").Return(config, nil)

	records := []schemas.CSVObjectFormat{
		{TotalIncome: 500000, WHT: 0, Donation: 0},
		{TotalIncome: 600000, WHT: 40000, Donation: 20000},
		{TotalIncome: 750000, WHT: 50000, Donation: 15000},
	}

	response, err := service.CalculateTaxFromCSV(records)
	assert.NoError(t, err)

	expectedTaxes := []schemas.CSVResponseMember{
		{TotalIncome: 500000, Tax: 29000},
		{TotalIncome: 600000, TaxRefund: 2000},
		{TotalIncome: 750000, Tax: 11250},
	}
	assert.Equal(t, expectedTaxes, response.Taxes)
}

func TestCalculateDetailedTax(t *testing.T) {
	mockRepo := new(MockTaxRepo)
	service := NewTaxService(mockRepo)
	config := &domains.TaxDeductionConfig{
		PersonalDeduction:    60000,
		DonationDeductionMax: 100000,
		KReceiptDeductionMax: 50000,
	}

	mockRepo.On("GetConfig").Return(config, nil)

	allowances := []schemas.Allowance{
		{AllowanceType: "k-receipt", Amount: 200000},
		{AllowanceType: "donation", Amount: 100000},
	}

	taxLevels, netTax, taxRefund, err := service.CalculateDetailedTax(900000, 7000, allowances)
	assert.NoError(t, err)

	expectedTaxLevels := []schemas.TaxLevel{
		{Level: "0-150,000", Tax: 0.0},
		{Level: "150,001-500,000", Tax: 35000.0},
		{Level: "500,001-1,000,000", Tax: 28500.0},
		{Level: "1,000,001-2,000,000", Tax: 0.0},
		{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
	}
	expectedNetTax := 56500.0
	expectedTaxRefund := 0.0

	assert.Equal(t, expectedTaxLevels, taxLevels)
	assert.Equal(t, expectedNetTax, netTax)
	assert.Equal(t, expectedTaxRefund, taxRefund)
}
