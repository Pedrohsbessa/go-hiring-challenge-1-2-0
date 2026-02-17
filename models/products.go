package models

import (
	"github.com/shopspring/decimal"
)

// Product represents a product stored in the catalog.
type Product struct {
	ID         uint            `gorm:"primaryKey"`
	Code       string          `gorm:"uniqueIndex;not null"`
	Price      decimal.Decimal `gorm:"type:decimal(10,2);not null"`
	CategoryID uint            `gorm:"not null"`
	Category   Category        `gorm:"foreignKey:CategoryID"`
	Variants   []Variant       `gorm:"foreignKey:ProductID"`
}

// TableName returns the database table name for Product.
func (p *Product) TableName() string {
	return "products"
}
