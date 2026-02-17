package models

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var ErrCategoryCodeAlreadyExists = errors.New("category code already exists")

type CategoriesRepository struct {
	db *gorm.DB
}

func NewCategoriesRepository(db *gorm.DB) *CategoriesRepository {
	return &CategoriesRepository{db: db}
}

func (r *CategoriesRepository) GetAllCategories() ([]Category, error) {
	var categories []Category
	if err := r.db.Order("id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoriesRepository) CreateCategory(category Category) (*Category, error) {
	if err := r.db.Create(&category).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrCategoryCodeAlreadyExists
		}

		return nil, err
	}

	return &category, nil
}
