package utilities

import (
	"fmt"
	"strings"

	"github.com/thitiphum-bluesage/assessment-tax/interfaces/schemas"
)

func ValidateUpdatePersonalDeductionRequest(req *schemas.UpdatePersonalDeductionRequest) error {
	if req.Amount == nil {
		return fmt.Errorf("amount is required")
	} else if *req.Amount < 10000 || *req.Amount > 100000 {
		return fmt.Errorf("amount must be between 10,000 and 100,000")
	}
	return nil
}

func ValidateUpdateKReceiptRequest(req *schemas.UpdateKReceiptRequest) error {
	if req.Amount == nil {
		return fmt.Errorf("amount for k-receipt is required")
	} else if *req.Amount < 1 || *req.Amount > 100000 {
		return fmt.Errorf("amount for k-receipt must be between 1 and 100,000")
	}
	return nil
}

func ValidateTaxCalculationRequest(req *schemas.TaxCalculationRequest) error {
	var errs []string

	// Check and dereference pointers for TotalIncome and WHT
	if req.TotalIncome == nil {
		errs = append(errs, "TotalIncome is required")
	} else if *req.TotalIncome < 0 {
		errs = append(errs, "TotalIncome must be non-negative")
	}

	if req.WHT == nil {
		errs = append(errs, "WHT is required")
	} else if *req.WHT < 0 {
		errs = append(errs, "WHT must be non-negative")
	} else if req.TotalIncome != nil && *req.WHT > *req.TotalIncome {
		errs = append(errs, "WHT cannot be greater than TotalIncome")
	}

	// Uncomment to not allow blank allowances
	// if len(req.Allowances) == 0 {
	// 	errs = append(errs, "At least one allowance must be provided")
	// }

	allowanceCounts := map[string]int{}
	for _, allowance := range req.Allowances {
		if allowance.AllowanceType != "donation" && allowance.AllowanceType != "k-receipt" {
			errs = append(errs, fmt.Sprintf("Invalid allowance type: %s. Allowed types are 'donation' and 'k-receipt'.", allowance.AllowanceType))
		}
		if allowance.Amount < 0 {
			errs = append(errs, fmt.Sprintf("Amount for %s must be non-negative", allowance.AllowanceType))
		}
		allowanceCounts[allowance.AllowanceType]++
		if allowanceCounts[allowance.AllowanceType] > 1 {
			errs = append(errs, fmt.Sprintf("Only one %s allowance can be included", allowance.AllowanceType))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errs, ", "))
	}
	return nil
}

func ValidateCSVTaxRecords(records []schemas.CSVObjectFormat) error {
	var errs []string

	for i, record := range records {
		if record.TotalIncome < 0 {
			errs = append(errs, fmt.Sprintf("Record %d: TotalIncome must be non-negative", i+1))
		}
		if record.WHT < 0 {
			errs = append(errs, fmt.Sprintf("Record %d: WHT must be non-negative", i+1))
		}
		if record.Donation < 0 {
			errs = append(errs, fmt.Sprintf("Record %d: Donation must be non-negative", i+1))
		}
		if record.KReceipt < 0 {
			errs = append(errs, fmt.Sprintf("Record %d: KReceipt must be non-negative", i+1))
		}
		if record.WHT > record.TotalIncome {
			errs = append(errs, fmt.Sprintf("Record %d: WHT cannot be greater than TotalIncome", i+1))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errs, "; "))
	}

	return nil
}
