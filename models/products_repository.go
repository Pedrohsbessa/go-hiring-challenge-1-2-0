package models

import (
	"strings"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

type ProductCatalogFilter struct {
	Offset        int
	Limit         int
	Category      string
	PriceLessThan *decimal.Decimal
}

func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) ListProducts(filter ProductCatalogFilter) ([]Product, int64, error) {
	query := r.db.Model(&Product{})

	if strings.TrimSpace(filter.Category) != "" {
		category := strings.TrimSpace(filter.Category)
		query = query.Joins("Category").Where("LOWER(categories.code) = LOWER(?) OR LOWER(categories.name) = LOWER(?)", category, category)
	}

	if filter.PriceLessThan != nil {
		query = query.Where("products.price < ?", *filter.PriceLessThan)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []Product
	if err := query.
		Preload("Category").
		Order("products.id ASC").
		Offset(filter.Offset).
		Limit(filter.Limit).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductsRepository) GetProductByCode(code string) (*Product, error) {
	var product Product
	if err := r.db.Preload("Category").Preload("Variants").Where("code = ?", code).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}
