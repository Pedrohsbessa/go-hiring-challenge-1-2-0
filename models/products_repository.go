package models

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ProductsRepository provides persistence operations for products.
type ProductsRepository struct {
	db *gorm.DB
}

// ProductCatalogFilter defines pagination and filter options for catalog listing.
type ProductCatalogFilter struct {
	Offset        int
	Limit         int
	Category      string
	PriceLessThan *decimal.Decimal
}

// NewProductsRepository creates a products repository backed by gorm.
func NewProductsRepository(db *gorm.DB) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

// ListProducts returns products and total count according to the provided filter.
func (r *ProductsRepository) ListProducts(filter ProductCatalogFilter) ([]Product, int64, error) {
	query := r.db.Model(&Product{})

	if strings.TrimSpace(filter.Category) != "" {
		category := strings.TrimSpace(filter.Category)
		query = query.Joins("Category").Where("LOWER(\"Category\".code) = LOWER(?) OR LOWER(\"Category\".name) = LOWER(?)", category, category)
	}

	if filter.PriceLessThan != nil {
		query = query.Where("products.price < ?", *filter.PriceLessThan)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count products failed: %w", err)
	}

	var products []Product
	if err := query.
		Preload("Category").
		Order("products.id ASC").
		Offset(filter.Offset).
		Limit(filter.Limit).
		Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("list products failed: %w", err)
	}

	return products, total, nil
}

// GetProductByCode returns a single product by code with category and variants preloaded.
func (r *ProductsRepository) GetProductByCode(code string) (*Product, error) {
	var product Product
	if err := r.db.Preload("Category").Preload("Variants").Where("code = ?", code).First(&product).Error; err != nil {
		return nil, fmt.Errorf("get product by code failed: %w", err)
	}

	return &product, nil
}
