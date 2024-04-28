package controllers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/thitiphum-bluesage/assessment-tax/applications/services/tax"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/schemas"
	"github.com/thitiphum-bluesage/assessment-tax/utilities"
)

type TaxController struct {
	taxService tax.TaxServiceInterface
}

func NewTaxController(service tax.TaxServiceInterface) *TaxController {
	return &TaxController{
		taxService: service,
	}
}

// used in story 1,2,3
func (tc *TaxController) CalculateTax(c echo.Context) error {
	var req schemas.TaxCalculationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input data")
	}

	if err := utilities.ValidateTaxCalculationRequest(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	netTax, taxRefund, err := tc.taxService.CalculateTax(*req.TotalIncome, *req.WHT, req.Allowances)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if taxRefund > 0 {
		response := schemas.TaxCalculationRefundResponse{
			TaxRefund: taxRefund,
		}
		return c.JSON(http.StatusOK, response)
	}

	response := schemas.TaxCalculationResponse{
		Tax: netTax,
	}
	return c.JSON(http.StatusOK, response)
}

// CalculateDetailedTax calculates the detailed tax amounts based on income, withholdings, and allowances
// @Summary Calculate detailed tax
// @Description Calculates taxes including breakdowns by tax level and potential refunds.
// @Tags tax
// @Accept json
// @Produce json
// @Param request body schemas.TaxCalculationRequest true "Tax Calculation Request"
// @Success 200 {object} schemas.DetailedTaxCalculationResponse "Detailed breakdown of tax calculations"
// @Failure 400 {object} schemas.ErrorResponse "Invalid input data"
// @Failure 500 {object} schemas.ErrorResponse "Internal Server Error"
// @Router /tax/calculations [post]
func (tc *TaxController) CalculateDetailedTax(c echo.Context) error {
	var req schemas.TaxCalculationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input data")
	}

	if err := utilities.ValidateTaxCalculationRequest(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	taxLevel, netTax, taxRefund, err := tc.taxService.CalculateDetailedTax(*req.TotalIncome, *req.WHT, req.Allowances)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if taxRefund > 0 {
		response := schemas.TaxCalculationRefundResponse{
			TaxRefund: taxRefund,
            TaxLevel: taxLevel,
		}
		return c.JSON(http.StatusOK, response)
	}

	response := schemas.DetailedTaxCalculationResponse{
		Tax:      netTax,
		TaxLevel: taxLevel,
	}
	return c.JSON(http.StatusOK, response)
}

// CalculateCSVTax calculates taxes from a CSV file upload containing multiple taxpayer records.
// @Summary Calculate taxes from CSV
// @Description Accepts a file upload (CSV format) with tax data, processes each record, and returns tax calculations.
// @Tags tax
// @Accept multipart/form-data
// @Produce json
// @Param taxFile formData file true "CSV file containing tax data"
// @Success 200 {object} schemas.CSVResponse "Tax calculations for all records in the uploaded CSV"
// @Failure 400 {object} schemas.ErrorResponse "Invalid input data or CSV format errors"
// @Failure 500 {object} schemas.ErrorResponse "Internal Server Error"
// @Router /tax/calculations/upload-csv [post]
func (tc *TaxController) CalculateCSVTax(c echo.Context) error {
    fileHeader, err := c.FormFile("taxFile")
    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Failed to get the file")
    }

    file, err := fileHeader.Open()
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to open the file")
    }
    defer file.Close()

    csvReader := csv.NewReader(file)
    headers, err := csvReader.Read()
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read headers from CSV file")
    }

    columnIndex := make(map[string]int)
    for i, header := range headers {
        columnIndex[strings.ToLower(header)] = i
    }

    hasTotalIncome := false
    for i, header := range headers {
        normalizedHeader := strings.ToLower(header)
        columnIndex[normalizedHeader] = i
        if normalizedHeader == "totalincome" {
            hasTotalIncome = true
        }
    }
    if !hasTotalIncome {
        return echo.NewHTTPError(http.StatusBadRequest, "CSV file does not contain 'totalincome' header")
    }

    // Validate expected columns
    for key:= range columnIndex {
        if key != "totalincome" && key != "wht" && key != "donation" && key != "k-receipt" {
            return echo.NewHTTPError(http.StatusBadRequest, "Invalid CSV file")
        }
    }

    var taxRecords []schemas.CSVObjectFormat
    for {
        record, err := csvReader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read record from CSV file")
        }

        var taxRecord schemas.CSVObjectFormat
        for key, index := range columnIndex {
            switch key {
            case "totalincome":
                if taxRecord.TotalIncome, err = parseAndValidateFloat(record[index]); err != nil {
                    return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid CSV file: %v", err))
                }
            case "wht":
                if taxRecord.WHT, err = parseAndValidateFloat(record[index]); err != nil {
                    return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid CSV file: %v", err))
                }
            case "donation":
                if taxRecord.Donation, err = parseAndValidateFloat(record[index]); err != nil {
                    return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid CSV file: %v", err))
                }
            case "k-receipt":
                if taxRecord.KReceipt, err = parseAndValidateFloat(record[index]); err != nil {
                    return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid CSV file: %v", err))
                }
            }
        }

        taxRecords = append(taxRecords, taxRecord)
    }

    if err := utilities.ValidateCSVTaxRecords(taxRecords); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    response, err := tc.taxService.CalculateTaxFromCSV(taxRecords)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }
    return c.JSON(http.StatusOK, response)
}

func parseAndValidateFloat(value string) (float64, error) {
    floatValue, err := strconv.ParseFloat(value, 64)
    if err != nil {
        return 0, err
    }
    return floatValue, nil
}