package catalog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestCatalogHandleGetByCodeSuccessWithVariantFallback(t *testing.T) {
	t.Parallel()

	variantPrice := decimal.RequireFromString("11.99")
	mock := &productsReaderMock{
		productByCode: &models.Product{
			Code:  "PROD001",
			Price: decimal.RequireFromString("10.99"),
			Category: models.Category{
				Code: "CLOTHING",
				Name: "Clothing",
			},
			Variants: []models.Variant{
				{Name: "Variant A", SKU: "SKU001A", Price: &variantPrice},
				{Name: "Variant B", SKU: "SKU001B", Price: nil},
			},
		},
	}

	handler := NewCatalogHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/catalog/PROD001", nil)
	req.SetPathValue("code", "PROD001")
	res := httptest.NewRecorder()

	handler.HandleGetByCode(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var payload ProductDetailsResponse
	err := json.Unmarshal(res.Body.Bytes(), &payload)
	assert.NoError(t, err)
	assert.Equal(t, "PROD001", payload.Code)
	assert.Equal(t, "CLOTHING", payload.Category.Code)
	assert.Len(t, payload.Variants, 2)
	assert.Equal(t, 11.99, payload.Variants[0].Price)
	assert.Equal(t, 10.99, payload.Variants[1].Price)
}

func TestCatalogHandleGetByCodeNotFound(t *testing.T) {
	t.Parallel()

	mock := &productsReaderMock{err: gorm.ErrRecordNotFound}
	handler := NewCatalogHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/catalog/MISSING", nil)
	req.SetPathValue("code", "MISSING")
	res := httptest.NewRecorder()

	handler.HandleGetByCode(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestCatalogHandleGetByCodeRepositoryError(t *testing.T) {
	t.Parallel()

	mock := &productsReaderMock{err: assert.AnError}
	handler := NewCatalogHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/catalog/PROD001", nil)
	req.SetPathValue("code", "PROD001")
	res := httptest.NewRecorder()

	handler.HandleGetByCode(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}
