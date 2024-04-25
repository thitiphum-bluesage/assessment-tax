package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thitiphum-bluesage/assessment-tax/applications/services/admin"
	"github.com/thitiphum-bluesage/assessment-tax/domains"
)

// Mock repository
type MockTaxDeductionConfigRepository struct {
	mock.Mock
}

func (m *MockTaxDeductionConfigRepository) GetConfig() (*domains.TaxDeductionConfig, error) {
	args := m.Called()
	if config, ok := args.Get(0).(*domains.TaxDeductionConfig); ok {
		return config, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTaxDeductionConfigRepository) UpdatePersonalDeduction(amount float64) error {
	args := m.Called(amount)
	return args.Error(0)
}

func (m *MockTaxDeductionConfigRepository) UpdateKReceiptDeductionMax(amount float64) error {
	args := m.Called(amount)
	return args.Error(0)
}

// Testing the AdminService with mocks
func TestAdminService_UpdatePersonalDeduction(t *testing.T) {
	mockRepo := new(MockTaxDeductionConfigRepository)
	adminService := admin.NewAdminService(mockRepo)

	// Test updating with a valid amount
	mockRepo.On("UpdatePersonalDeduction", 70000.0).Return(nil)
	err := adminService.UpdatePersonalDeduction(70000.0)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Test updating with an invalid amount (too low)
	err = adminService.UpdatePersonalDeduction(9000)
	assert.Error(t, err, "amount must be between 10,000 and 100,000")

	// Test updating with an invalid amount (too high)
	err = adminService.UpdatePersonalDeduction(101000)
	assert.Error(t, err, "amount must be between 10,000 and 100,000")
}

func TestAdminService_UpdateKReceiptDeductionMax(t *testing.T) {
	mockRepo := new(MockTaxDeductionConfigRepository)
	adminService := admin.NewAdminService(mockRepo)

	// Test updating within valid range
	mockRepo.On("UpdateKReceiptDeductionMax", 50000.0).Return(nil)
	err := adminService.UpdateKReceiptDeductionMax(50000.0)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Test updating with a negative amount
	err = adminService.UpdateKReceiptDeductionMax(-100.0)
	assert.Error(t, err, "amount must be less than or equal to 100,000")

	// Test updating with an amount too high
	err = adminService.UpdateKReceiptDeductionMax(100001.0)
	assert.Error(t, err, "amount must be less than or equal to 100,000")
}
