package models

import "github.com/google/uuid"

type CartItem struct {
	ID        uint         `gorm:"primaryKey"`
	CartID    uuid.UUID    `gorm:"type:uuid;not null" json:"cart_id"`
	Cart      ShoppingCart `gorm:"foreignKey:CartID;references:UUID;constraint:OnDelete:CASCADE"`
	ProductID uint         `gorm:"not null"`
	Product   Product      `gorm:"foreignKey:ProductID"`
	Quantity  int          `gorm:"not null"`
}
type WishList struct {
	ID        uint    `gorm:"primaryKey"`
	ProductID uint    `gorm:"not null" json:"-"`
	Product   Product `gorm:"foreignKey:ProductID"`
	UserID    uint    `gorm:"not null" json:"-"`
	User      User    `gorm:"foreignKey:UserID" json:"-"`
}
