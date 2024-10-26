package models

type OrderItem struct {
	ID              uint    `gorm:"primaryKey"`
	OrderID         uint    `gorm:"not null"`
	Order           Order   `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	ProductID       uint    `gorm:"not null"`
	Product         Product `gorm:"foreignKey:ProductID"`
	Quantity        int     `gorm:"not null"`
	PriceAtPurchase float64 `gorm:"type:decimal(10,2);not null"`
}
