package categories

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/stretchr/testify/assert"
)

type categoriesRepoMock struct {
	categories       []models.Category
	getErr           error
	createErr        error
	capturedCategory *models.Category
}

func (m *categoriesRepoMock) GetAllCategories() ([]models.Category, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}

	return m.categories, nil
}

func (m *categoriesRepoMock) CreateCategory(category models.Category) (*models.Category, error) {
	m.capturedCategory = &category
	if m.createErr != nil {
		return nil, m.createErr
	}

	return &category, nil
}

func TestHandleGetCategoriesSuccess(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&categoriesRepoMock{
		categories: []models.Category{
			{Code: "CLOTHING", Name: "Clothing"},
			{Code: "SHOES", Name: "Shoes"},
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	res := httptest.NewRecorder()

	handler.HandleGet(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	var payload ListResponse
	err := json.Unmarshal(res.Body.Bytes(), &payload)
	assert.NoError(t, err)
	assert.Len(t, payload.Categories, 2)
	assert.Equal(t, "CLOTHING", payload.Categories[0].Code)
}

func TestHandlePostCategorySuccess(t *testing.T) {
	t.Parallel()

	mock := &categoriesRepoMock{}
	handler := NewHandler(mock)

	body := []byte(`{"code":"BAGS","name":"Bags"}`)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	res := httptest.NewRecorder()

	handler.HandlePost(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.NotNil(t, mock.capturedCategory)
	assert.Equal(t, "BAGS", mock.capturedCategory.Code)
	assert.Equal(t, "Bags", mock.capturedCategory.Name)
}

func TestHandlePostCategoryValidationError(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&categoriesRepoMock{})

	body := []byte(`{"code":"","name":""}`)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	res := httptest.NewRecorder()

	handler.HandlePost(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestHandlePostCategoryConflict(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&categoriesRepoMock{createErr: errors.New("duplicate key value violates unique constraint")})

	body := []byte(`{"code":"BAGS","name":"Bags"}`)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	res := httptest.NewRecorder()

	handler.HandlePost(res, req)

	assert.Equal(t, http.StatusConflict, res.Code)
}
