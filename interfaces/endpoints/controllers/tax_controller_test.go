package controllers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/schemas"
)

type MockTaxService struct {
	mock.Mock
}

// Mock implementation of CalculateTax
func (m *MockTaxService) CalculateTax(totalIncome float64, wht float64, allowances []schemas.Allowance) (float64, float64, error) {
	args := m.Called(totalIncome, wht, allowances)
	return args.Get(0).(float64), args.Get(1).(float64), args.Error(2)
}

// Mock implementation of CalculateDetailedTax
func (m *MockTaxService) CalculateDetailedTax(totalIncome, wht float64, allowances []schemas.Allowance) ([]schemas.TaxLevel, float64, float64, error) {
	args := m.Called(totalIncome, wht, allowances)
	return args.Get(0).([]schemas.TaxLevel), args.Get(1).(float64), args.Get(2).(float64), args.Error(3)
}

// Mock implementation of CalculateTaxFromCSV
func (m *MockTaxService) CalculateTaxFromCSV(records []schemas.CSVObjectFormat) (schemas.CSVResponse, error) {
	args := m.Called(records)
	return args.Get(0).(schemas.CSVResponse), args.Error(1)
}

func TestTaxController_CalculateDetailedTax_Success(t *testing.T) {
	e := echo.New()
	mockService := new(MockTaxService)
	controller := NewTaxController(mockService)

	// Setting up the mock response
	taxLevels := []schemas.TaxLevel{{Level: "Basic", Tax: 5000}}
	mockService.On("CalculateDetailedTax", 100000.0, 10000.0, []schemas.Allowance{}).Return(taxLevels, 90000.0, 0.0, nil)

	reqBody := `{"TotalIncome": 100000, "WHT": 10000, "Allowances": []}`
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test the controller method
	if assert.NoError(t, controller.CalculateDetailedTax(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp schemas.DetailedTaxCalculationResponse
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp)) {
			assert.Equal(t, 90000.0, resp.Tax)
			assert.Len(t, resp.TaxLevel, 1)
		}
	}
}

func TestTaxController_CalculateDetailedTax_InvalidFile_WHTGreaterThanTotalIncome(t *testing.T) {
    e := echo.New()
    controller := NewTaxController(nil)

    reqBody := `{
        "totalIncome": 900000,
        "wht": 1000000,
        "allowances": [
            {
                "allowanceType": "k-receipt",
                "amount": 200000.0
            },
            {
                "allowanceType": "donation",
                "amount": 100000.0
            }
        ]
    }`
    req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(reqBody))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    err := controller.CalculateDetailedTax(c)
    assert.Error(t, err)
    assert.IsType(t, &echo.HTTPError{}, err)
    assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func TestTaxController_CalculateDetailedTax_InvalidFile_NegativeTotalIncome(t *testing.T) {
    e := echo.New()
    controller := NewTaxController(nil)

    reqBody := `{
        "totalIncome": -900000,
        "wht": 7000,
        "allowances": [
            {
                "allowanceType": "k-receipt",
                "amount": 200000.0
            },
            {
                "allowanceType": "donation",
                "amount": 100000.0
            }
        ]
    }`
    req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(reqBody))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    err := controller.CalculateDetailedTax(c)
    assert.Error(t, err)
    assert.IsType(t, &echo.HTTPError{}, err)
    assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func TestTaxController_CalculateDetailedTax_InvalidFile_NegativeAllowanceAmount(t *testing.T) {
    e := echo.New()
    controller := NewTaxController(nil)

    reqBody := `{
        "totalIncome": 900000,
        "wht": 7000,
        "allowances": [
            {
                "allowanceType": "k-receipt",
                "amount": -200000.0
            },
            {
                "allowanceType": "donation",
                "amount": 100000.0
            }
        ]
    }`
    req := httptest.NewRequest(http.MethodPost, "/tax/calculations", strings.NewReader(reqBody))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    err := controller.CalculateDetailedTax(c)
    assert.Error(t, err)
    assert.IsType(t, &echo.HTTPError{}, err)
    assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}


func TestTaxController_CalculateCSVTax_ValidFile(t *testing.T) {
    e := echo.New()

    // Create a test CSV file
    csvData := "totalIncome,wht,donation\n500000,0,0\n600000,40000,20000\n750000,50000,15000\n"
    body := new(bytes.Buffer)
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile("taxFile", "test.csv")
    assert.NoError(t, err)
    _, err = part.Write([]byte(csvData))
    assert.NoError(t, err)
    err = writer.Close()
    assert.NoError(t, err)

    // Create a new request
    req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    rec := httptest.NewRecorder()

    // Create a new MockTaxService instance
    mockTaxService := new(MockTaxService)

    // Set up the mock expectations
    expectedTaxRecords := []schemas.CSVObjectFormat{
        {TotalIncome: 500000, WHT: 0, Donation: 0},
        {TotalIncome: 600000, WHT: 40000, Donation: 20000},
        {TotalIncome: 750000, WHT: 50000, Donation: 15000},
    }
    expectedResponse := schemas.CSVResponse{
        Taxes: []schemas.CSVResponseMember{
            {TotalIncome: 500000, Tax: 29000},
            {TotalIncome: 600000, TaxRefund: 2000},
            {TotalIncome: 750000, Tax: 11250},
        },
    }
    mockTaxService.On("CalculateTaxFromCSV", expectedTaxRecords).Return(expectedResponse, nil)

    taxController := &TaxController{
        taxService: mockTaxService,
    }

    if assert.NoError(t, taxController.CalculateCSVTax(e.NewContext(req, rec))) {
        assert.Equal(t, http.StatusOK, rec.Code)
        // Add assertions for the response body
        var resp schemas.CSVResponse
        if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp)) {
            assert.Equal(t, expectedResponse.Taxes, resp.Taxes)
        }
    }

    mockTaxService.AssertExpectations(t)
}

func TestTaxController_CalculateCSVTax_InvalidFile_WHTGreaterThanTotalIncome(t *testing.T) {
    // Create a new Echo instance
    e := echo.New()

    // Create a test CSV file with WHT greater than TotalIncome
    csvData := "totalIncome,wht,donation\n500000,600000,0\n"
    body := new(bytes.Buffer)
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile("taxFile", "test.csv")
    assert.NoError(t, err)
    _, err = part.Write([]byte(csvData))
    assert.NoError(t, err)
    err = writer.Close()
    assert.NoError(t, err)

    // Create a new request
    req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    rec := httptest.NewRecorder()

    // Create a new TaxController instance
    taxController := &TaxController{
        taxService: &MockTaxService{},
    }

    err = taxController.CalculateCSVTax(e.NewContext(req, rec))
    assert.Error(t, err)
    assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func TestTaxController_CalculateCSVTax_InvalidFile_NegativeValues(t *testing.T) {
    // Create a new Echo instance
    e := echo.New()

    // Create a test CSV file with negative values
    csvData := "totalIncome,wht,donation\n-500000,0,0\n"
    body := new(bytes.Buffer)
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile("taxFile", "test.csv")
    assert.NoError(t, err)
    _, err = part.Write([]byte(csvData))
    assert.NoError(t, err)
    err = writer.Close()
    assert.NoError(t, err)

    // Create a new request
    req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    rec := httptest.NewRecorder()

    // Create a new TaxController instance
    taxController := &TaxController{
        taxService: &MockTaxService{},
    }

    err = taxController.CalculateCSVTax(e.NewContext(req, rec))
    assert.Error(t, err)
    assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func TestTaxController_CalculateCSVTax_InvalidFile_StringValues(t *testing.T) {
    // Create a new Echo instance
    e := echo.New()

    // Create a test CSV file with string values
    csvData := "totalIncome,wht,donation\nGodOuIsHere,0,0\n"
    body := new(bytes.Buffer)
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile("taxFile", "test.csv")
    assert.NoError(t, err)
    _, err = part.Write([]byte(csvData))
    assert.NoError(t, err)
    err = writer.Close()
    assert.NoError(t, err)

    // Create a new request
    req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    rec := httptest.NewRecorder()

    // Create a new TaxController instance
    taxController := &TaxController{
        taxService: &MockTaxService{},
    }

    err = taxController.CalculateCSVTax(e.NewContext(req, rec))
    assert.Error(t, err)
    assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}
