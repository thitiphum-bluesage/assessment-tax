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

type TaxCalculationRefundResponse struct {
	TaxRefund float64 `json:"taxRefund"`
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type DetailedTaxCalculationResponse struct {
	Tax      float64    `json:"tax"`
	TaxLevel []TaxLevel `json:"taxLevel"`
}

type CSVObjectFormat struct {
	TotalIncome float64 `csv:"totalIncome"`
	WHT         float64 `csv:"wht"`
	Donation    float64 `csv:"donation"`
	KReceipt    float64 `csv:"k-receipt"`
}

type CSVResponseMember struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax,omitempty"`
	TaxRefund   float64 `json:"taxRefund,omitempty"`
}

type CSVResponse struct {
	Taxes []CSVResponseMember `json:"taxes"`
}
