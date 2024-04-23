package schemas

type UpdatePersonalDeductionRequest struct {
	Amount *float64 `json:"amount"`
}

type UpdatePersonalDeductionResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type UpdateKReceiptRequest struct {
	Amount *float64 `json:"amount" `
}

type UpdateKReceiptResponse struct {
	KReceipt float64 `json:"kReceipt"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType" `
	Amount        float64 `json:"amount" `
}

type TaxCalculationRequest struct {
	TotalIncome *float64    `json:"totalIncome" `
	WHT         *float64    `json:"wht" `
	Allowances  []Allowance `json:"allowances" `
}

type TaxCalculationResponse struct {
	Tax float64 `json:"tax"`
}
