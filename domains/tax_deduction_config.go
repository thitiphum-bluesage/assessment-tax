package domains

type TaxDeductionConfig struct {
	ConfigName           string  `gorm:"type:varchar(100);not null;unique"`
	PersonalDeduction    float64 `gorm:"type:float;not null;check:personal_deduction >= 10000 and personal_deduction <= 100000"`
	KReceiptDeductionMax float64 `gorm:"type:float;not null;check:k_receipt_deduction_max <= 100000"`
	DonationDeductionMax float64 `gorm:"type:float;not null;default:100000"`
}
