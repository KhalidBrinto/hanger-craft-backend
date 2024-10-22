package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string         `gorm:"size:150;not null" json:"name"`
	Description string         `json:"description"`
	SKU         string         `gorm:"size:150;not null;unique;index" json:"sku"`
	Barcode     *string        `gorm:"size:150" json:"barcode"`
	Price       float64        `gorm:"not null" json:"price"`
	Currency    string         `gorm:"size:3; not null" json:"currency"`
	Images      pq.StringArray `json:"images"`
	CategoryID  uint           `gorm:"not null" json:"category_id"`
	Category    Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Reviews     []Review       `gorm:"foreignKey:ProductID" json:"reviews"`
}

type ProductAttribute struct {
	gorm.Model
	Name        string  `gorm:"size:150;not null" json:"name"`
	Description string  `json:"description"`
	ProductID   uint    `gorm:"not null" json:"product_id"`
	Product     Product `gorm:"foreignKey:ProductID" json:"product"`
}
