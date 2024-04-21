package repository

import (
	"github.com/thitiphum-bluesage/assessment-tax/domains"
	"gorm.io/gorm"
)

type taxDeductionConfigRepository struct {
    db *gorm.DB
}

func NewTaxDeductionConfigRepository(db *gorm.DB) TaxDeductionConfigRepository {
    return &taxDeductionConfigRepository{db: db}
}

func (r *taxDeductionConfigRepository) GetConfig() (*domains.TaxDeductionConfig, error) {
    var config domains.TaxDeductionConfig
    err := r.db.Where("config_name = ?", "MainConfig").First(&config).Error
    if err != nil {
        return nil, err
    }
    return &config, nil
}

func (r *taxDeductionConfigRepository) UpdatePersonalDeduction(amount float64) error {
    result := r.db.Model(&domains.TaxDeductionConfig{}).Where("config_name = ?", "MainConfig").Update("personal_deduction", amount)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func (r *taxDeductionConfigRepository) UpdateKReceiptDeductionMax(amount float64) error {
    result := r.db.Model(&domains.TaxDeductionConfig{}).Where("config_name = ?", "MainConfig").Update("k_receipt_deduction_max", amount)
    if result.Error != nil {
        return result.Error
    }
    return nil
}
