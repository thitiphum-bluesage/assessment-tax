package controllers

import (
	"encoding/csv"
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
		}
		return c.JSON(http.StatusOK, response)
	}

	response := schemas.DetailedTaxCalculationResponse{
		Tax:      netTax,
		TaxLevel: taxLevel,
	}
	return c.JSON(http.StatusOK, response)
}

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

	// {
	// 	"totalincome": 0,
	// 	"wht": 1,
	// 	"donation": 2,
	// 	"kreceipt": 3
	// }

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
		if index, ok := columnIndex["totalincome"]; ok {
			taxRecord.TotalIncome, _ = strconv.ParseFloat(record[index], 64)
		}
		if index, ok := columnIndex["wht"]; ok {
			taxRecord.WHT, _ = strconv.ParseFloat(record[index], 64)
		}
		if index, ok := columnIndex["donation"]; ok {
			taxRecord.Donation, _ = strconv.ParseFloat(record[index], 64)
		}
		// k-receipt will = 0 if not provided
		if index, ok := columnIndex["k-receipt"]; ok {
			taxRecord.KReceipt, _ = strconv.ParseFloat(record[index], 64)
		}

		taxRecords = append(taxRecords, taxRecord)
	}

	response, err := tc.taxService.CalculateTaxFromCSV(taxRecords)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, response)
}
