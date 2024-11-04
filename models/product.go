package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string   `gorm:"size:150;not null"`
	Description string   `gorm:"type:text"`
	SKU         string   `gorm:"size:150;not null;unique;index"`
	Barcode     *string  `gorm:"size:150"`
	Price       float64  `gorm:"type:decimal(10,2);not null"`
	Currency    string   `gorm:"size:3; not null"`
	CategoryID  uint     `gorm:"not null"`
	Category    Category `gorm:"foreignKey:CategoryID"`
	Status      *string  `gorm:"not null;check:status IN ('published', 'unpublished')"`
	Featured    bool     `gorm:"default:false"`
	Stock       uint     `gorm:"-"`
	IsChild     bool     `gorm:"default:false"`
	ParentID    *uint
	Color       string
	Size        string
	BrandID     *uint
	Brand       Brand          `gorm:"foreignKey:BrandID"`
	Images      []ProductImage `gorm:"foreignKey:ProductID"`
}

type ProductImage struct {
	ID        uint    `gorm:"primaryKey"`
	ProductID uint    // Foreign key to the Product
	Product   Product `gorm:"foreignKey:ProductID" json:"-"`
	ImageB64  string  `gorm:"-" json:"Image"`
	Image     []byte  `gorm:"type:bytea"` // Binary data for the image
}

func (c *ProductImage) BeforeCreate(tx *gorm.DB) (err error) {

	return nil

}

type ProductAttribute struct {
	gorm.Model
	Name        string  `gorm:"size:150;not null"`
	Description string  `gorm:"type:text"`
	ProductID   uint    `gorm:"not null"`
	Product     Product `gorm:"foreignKey:ProductID"`
}
