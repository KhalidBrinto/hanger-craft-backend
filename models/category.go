package models

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name         string `gorm:"size:100;not null"`
	CategoryType string `gorm:"size:100;not null"`
	ParentID     *uint
	Products     []Product `gorm:"foreignKey:CategoryID"`
}
