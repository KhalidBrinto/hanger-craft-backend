package models

import (
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	UserID    uint    `gorm:"not null"`
	User      User    `gorm:"foreignKey:UserID"`
	ProductID uint    `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID" json:"-"`
	Rating    int     `gorm:"check:rating >= 1 AND rating <= 5"`
	Comment   string  `gorm:"type:text"`
}
