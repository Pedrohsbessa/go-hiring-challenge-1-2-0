package catalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type productsReaderMock struct {
	products      []models.Product
	total         int64
	err           error
	capturedQuery models.ProductCatalogFilter
	productByCode *models.Product
}

func (m *productsReaderMock) ListProducts(filter models.ProductCatalogFilter) ([]models.Product, int64, error) {
	m.capturedQuery = filter
	return m.products, m.total, m.err
}

func (m *productsReaderMock) GetProductByCode(code string) (*models.Product, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.productByCode, nil
}

func TestCatalogHandleGetDefaults(t *testing.T) {
	t.Parallel()

	mock := &productsReaderMock{
		products: []models.Product{
			{
				Code:  "PROD001",
				Price: decimal.RequireFromString("10.99"),
				Category: models.Category{
					Code: "CLOTHING",
					Name: "Clothing",
				},
			},
		},
		total: 8,
	}

	handler := NewCatalogHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
	res := httptest.NewRecorder()

	handler.HandleGet(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, 0, mock.capturedQuery.Offset)
	assert.Equal(t, 10, mock.capturedQuery.Limit)
	assert.Equal(t, "", mock.capturedQuery.Category)
	assert.Nil(t, mock.capturedQuery.PriceLessThan)

	var payload Response
	err := json.Unmarshal(res.Body.Bytes(), &payload)
	assert.NoError(t, err)
	assert.Equal(t, int64(8), payload.Total)
	assert.Len(t, payload.Products, 1)
	assert.Equal(t, "PROD001", payload.Products[0].Code)
	assert.Equal(t, "CLOTHING", payload.Products[0].Category.Code)
	assert.Equal(t, "Clothing", payload.Products[0].Category.Name)
}

func TestCatalogHandleGetFiltersAndPagination(t *testing.T) {
	t.Parallel()

	mock := &productsReaderMock{}
	handler := NewCatalogHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/catalog?offset=3&limit=250&category=Shoes&price_lt=12.50", nil)
	res := httptest.NewRecorder()

	handler.HandleGet(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, 3, mock.capturedQuery.Offset)
	assert.Equal(t, 100, mock.capturedQuery.Limit)
	assert.Equal(t, "Shoes", mock.capturedQuery.Category)
	assert.NotNil(t, mock.capturedQuery.PriceLessThan)
	assert.True(t, decimal.RequireFromString("12.50").Equal(*mock.capturedQuery.PriceLessThan))
}

func TestCatalogHandleGetInvalidPriceLt(t *testing.T) {
	t.Parallel()

	mock := &productsReaderMock{}
	handler := NewCatalogHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/catalog?price_lt=invalid", nil)
	res := httptest.NewRecorder()

	handler.HandleGet(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestCatalogHandleGetRepositoryError(t *testing.T) {
	t.Parallel()

	mock := &productsReaderMock{err: errors.New("db failed")}
	handler := NewCatalogHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/catalog", nil)
	res := httptest.NewRecorder()

	handler.HandleGet(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
}
