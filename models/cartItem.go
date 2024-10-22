package models

type CartItem struct {
	ID        uint         `gorm:"primaryKey"`
	CartID    uint         `gorm:"not null" json:"cart_id"`
	Cart      ShoppingCart `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
	ProductID uint         `gorm:"not null"`
	Product   Product      `gorm:"foreignKey:ProductID"`
	Quantity  int          `gorm:"not null"`
}
