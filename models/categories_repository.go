package models

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// ErrCategoryCodeAlreadyExists indicates a unique violation for category code.
var ErrCategoryCodeAlreadyExists = errors.New("category code already exists")

// CategoriesRepository provides persistence operations for categories.
type CategoriesRepository struct {
	db *gorm.DB
}

// NewCategoriesRepository creates a categories repository backed by gorm.
func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{db: db}
}

// GetAllCategories returns all categories ordered by id.
func (r *CategoriesRepository) GetAllCategories() ([]Category, error) {
	var categories []Category
	if err := r.db.Order("id ASC").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("list categories failed: %w", err)
	}

	return categories, nil
}

// CreateCategory persists a new category.
func (r *CategoriesRepository) CreateCategory(category Category) (*Category, error) {
	if err := r.db.Create(&category).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrCategoryCodeAlreadyExists
		}

		return nil, fmt.Errorf("create category failed: %w", err)
	}

	return &category, nil
}
