package repository

import "github.com/thitiphum-bluesage/assessment-tax/domains"

type TaxDeductionConfigRepositoryInterface interface {
	GetConfig() (*domains.TaxDeductionConfig, error)
	UpdatePersonalDeduction(amount float64) error
	UpdateKReceiptDeductionMax(amount float64) error
}
