package models

import (
	"github.com/shopspring/decimal"
)

// Variant represents a product variant.
type Variant struct {
	ID        uint             `gorm:"primaryKey"`
	ProductID uint             `gorm:"not null"`
	Name      string           `gorm:"not null"`
	SKU       string           `gorm:"uniqueIndex;not null"`
	Price     *decimal.Decimal `gorm:"type:decimal(10,2);null"`
}

// TableName returns the database table name for Variant.
func (v *Variant) TableName() string {
	return "product_variants"
}
