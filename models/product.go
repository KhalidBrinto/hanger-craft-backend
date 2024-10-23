package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string `gorm:"size:150;not null"`
	Description string
	SKU         string         `gorm:"size:150;not null;unique;index"`
	Barcode     *string        `gorm:"size:150"`
	Price       float64        `gorm:"not null"`
	Currency    string         `gorm:"size:3; not null"`
	Images      pq.StringArray `gorm:"type:varchar[]"`
	CategoryID  uint           `gorm:"not null"`
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
