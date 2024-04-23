package services

import (
	"errors"

	"github.com/thitiphum-bluesage/assessment-tax/infrastructure/repository"
)

type adminService struct {
	taxRepo repository.TaxDeductionConfigRepositoryInterface
}

func NewAdminService(taxRepo repository.TaxDeductionConfigRepositoryInterface) AdminServiceInterface {
	return &adminService{
		taxRepo: taxRepo,
	}
}

func (s *adminService) UpdatePersonalDeduction(amount float64) error {
	if amount < 10000 || amount > 100000 {
		return errors.New("amount must be between 10,000 and 100,000")
	}
	return s.taxRepo.UpdatePersonalDeduction(amount)
}

func (s *adminService) UpdateKReceiptDeductionMax(amount float64) error {
	if amount < 0 || amount > 100000 {
		return errors.New("amount must be less than or equal to 100,000")
	}
	return s.taxRepo.UpdateKReceiptDeductionMax(amount)
}
