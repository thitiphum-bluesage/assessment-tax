package admin

type AdminServiceInterface interface {
	UpdatePersonalDeduction(amount float64) error
	UpdateKReceiptDeductionMax(amount float64) error
}
