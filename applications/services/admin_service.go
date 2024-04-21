package services

import (
	"errors"

	"github.com/thitiphum-bluesage/assessment-tax/infrastructure/repository"
)

type AdminService struct {
    taxRepo repository.TaxDeductionConfigRepository
}

func NewAdminService(taxRepo repository.TaxDeductionConfigRepository) *AdminService {
    return &AdminService{
        taxRepo: taxRepo,
    }
}

// UpdatePersonalDeduction updates the personal deduction amount.
func (s *AdminService) UpdatePersonalDeduction(amount float64) error {
    return s.taxRepo.UpdatePersonalDeduction(amount)
}

// UpdateKReceiptDeductionMax updates the maximum deduction for K-receipts.
func (s *AdminService) UpdateKReceiptDeductionMax(amount float64) error {
    if amount < 0 || amount > 100000 {
        return errors.New("amount must be less than or equal to 100,000")
    }
    return s.taxRepo.UpdateKReceiptDeductionMax(amount)
}
