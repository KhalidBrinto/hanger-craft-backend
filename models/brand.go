package models

import (
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type Brand struct {
	gorm.Model
	Name     null.String `gorm:"size:100;not null"`
	Products []Product   `gorm:"foreignKey:BrandID"`
}

// func (c *Brand) BeforeCreate(tx *gorm.DB) (err error) {

// 	if c.CategoryType.String == "child" && c.ParentID == nil {
// 		return errors.New("must provide parent category id when creating child category")
// 	}

// 	return nil

// }
