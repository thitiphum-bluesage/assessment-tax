package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thitiphum-bluesage/assessment-tax/interfaces/schemas"
)

type MockAdminService struct {
	mock.Mock
}

func (m *MockAdminService) UpdatePersonalDeduction(amount float64) error {
	args := m.Called(amount)
	return args.Error(0)
}

func (m *MockAdminService) UpdateKReceiptDeductionMax(amount float64) error {
	args := m.Called(amount)
	return args.Error(0)
}

func TestAdminController_UpdatePersonalDeduction_ValidInput(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new request with valid input
	validAmount := 50000.0
	reqBody := schemas.UpdatePersonalDeductionRequest{Amount: &validAmount}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/admin/personal-deduction", strings.NewReader(string(jsonBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a new MockAdminService instance
	mockService := new(MockAdminService)

	// Set up the mock expectation
	mockService.On("UpdatePersonalDeduction", validAmount).Return(nil)

	// Create a new AdminController instance with the mock service
	controller := &AdminController{
		service: mockService,
	}

	if assert.NoError(t, controller.UpdatePersonalDeduction(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp schemas.UpdatePersonalDeductionResponse
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp)) {
			assert.Equal(t, validAmount, resp.PersonalDeduction)
		}
	}

	mockService.AssertExpectations(t)
}

func TestAdminController_UpdatePersonalDeduction_InvalidInput(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new request with invalid input
	invalidAmount := 5000.0
	reqBody := schemas.UpdatePersonalDeductionRequest{Amount: &invalidAmount}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/admin/personal-deduction", strings.NewReader(string(jsonBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a new AdminController instance
	controller := &AdminController{
		service: nil,
	}

	err := controller.UpdatePersonalDeduction(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "amount must be between 10,000 and 100,000")
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func TestAdminController_UpdateKReceiptDeduction_ValidInput(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new request with valid input
	validAmount := 50000.0
	reqBody := schemas.UpdateKReceiptRequest{Amount: &validAmount}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/admin/k-receipt-deduction", strings.NewReader(string(jsonBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a new MockAdminService instance
	mockService := new(MockAdminService)

	// Set up the mock expectation
	mockService.On("UpdateKReceiptDeductionMax", validAmount).Return(nil)

	// Create a new AdminController instance with the mock service
	controller := &AdminController{
		service: mockService,
	}

	if assert.NoError(t, controller.UpdateKReceiptDeduction(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp schemas.UpdateKReceiptResponse
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp)) {
			assert.Equal(t, validAmount, resp.KReceipt)
		}
	}

	mockService.AssertExpectations(t)
}

func TestAdminController_UpdateKReceiptDeduction_InvalidInput(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	// Create a new request with invalid input
	invalidAmount := 0.0
	reqBody := schemas.UpdateKReceiptRequest{Amount: &invalidAmount}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/admin/k-receipt-deduction", strings.NewReader(string(jsonBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Create a new AdminController instance
	controller := &AdminController{
		service: nil,
	}

	err := controller.UpdateKReceiptDeduction(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "amount for k-receipt must be between 1 and 100,000")
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}
