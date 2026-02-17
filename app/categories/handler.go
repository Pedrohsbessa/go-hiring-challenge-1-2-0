package categories

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
)

type CategoryReaderWriter interface {
	GetAllCategories() ([]models.Category, error)
	CreateCategory(category models.Category) (*models.Category, error)
}

type Handler struct {
	repo CategoryReaderWriter
}

func NewHandler(repo CategoryReaderWriter) *Handler {
	return &Handler{repo: repo}
}

type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ListResponse struct {
	Categories []CategoryResponse `json:"categories"`
}

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

type CreateCategoryRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

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
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			api.ErrorResponse(w, http.StatusConflict, "category code already exists")
			return
		}

		api.ErrorResponse(w, http.StatusInternalServerError, "failed to create category")
		return
	}

	api.OKResponse(w, CategoryResponse{Code: created.Code, Name: created.Name})
}
