package models

import (
	"github.com/google/uuid"
)

type ShoppingCart struct {
	UUID      uuid.UUID  `gorm:"type:uuid;primaryKey;index"`
	UserID    uint       `gorm:"not null"`
	User      User       `gorm:"foreignKey:UserID"`
	CartItems []CartItem `gorm:"foreignKey:CartID"`
	// CartItems []CartItem `gorm:"foreignKey:CartID"`
}
