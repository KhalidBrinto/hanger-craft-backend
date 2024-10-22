package models

import "gorm.io/gorm"

type ShoppingCart struct {
	gorm.Model
	UserID    uint       `gorm:"not null"`
	User      User       `gorm:"foreignKey:UserID"`
	CartItems []CartItem `gorm:"foreignKey:CartID"`
}
