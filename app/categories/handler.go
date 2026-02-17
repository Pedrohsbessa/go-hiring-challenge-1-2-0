package categories

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

// CategoryReaderWriter defines category operations consumed by the handler.
type CategoryReaderWriter interface {
	GetAllCategories() ([]models.Category, error)
	CreateCategory(category models.Category) (*models.Category, error)
}

// Handler exposes HTTP handlers for category endpoints.
type Handler struct {
	repo CategoryReaderWriter
}

// NewHandler creates a new category handler.
func NewHandler(repo CategoryReaderWriter) *Handler {
	return &Handler{repo: repo}
}

// CategoryResponse represents category data returned by API responses.
type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// ListResponse contains category list payload.
type ListResponse struct {
	Categories []CategoryResponse `json:"categories"`
}

// HandleGet returns all categories.
func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repo.GetAllCategories()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "failed to fetch categories")
		return
	}

	response := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		response[i] = CategoryResponse{Code: category.Code, Name: category.Name}
	}

	api.OKResponse(w, ListResponse{Categories: response})
}

// CreateCategoryRequest represents category creation payload.
type CreateCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// HandlePost validates and creates a new category.
func (h *Handler) HandlePost(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)
	if req.Code == "" || req.Name == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "code and name are required")
		return
	}

	created, err := h.repo.CreateCategory(models.Category{Code: req.Code, Name: req.Name})
	if err != nil {
		if errors.Is(err, models.ErrCategoryCodeAlreadyExists) {
			api.ErrorResponse(w, http.StatusConflict, "category code already exists")
			return
		}

		api.ErrorResponse(w, http.StatusInternalServerError, "failed to create category")
		return
	}

	api.CreatedResponse(w, CategoryResponse{Code: created.Code, Name: created.Name})
}
