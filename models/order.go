package models

import (
	"backend/utils"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	OrderIdentifier      string           `gorm:"type:varchar(8); not null;unique;index"`
	UserID               uint             `gorm:"not null"`
	User                 User             `gorm:"foreignKey:UserID"`
	OrderStatus          string           `gorm:"size:50;not null;check:order_status IN ('pending', 'shipped', 'delivered', 'cancelled')"`
	Currency             *string          `gorm:"size:3; not null"`
	TotalPrice           float64          `gorm:"type:decimal(10,2);not null"`
	ItemPrice            float64          `gorm:"type:decimal(10,2);not null"`
	DiscountAmount       float64          `gorm:"type:decimal(10,2);default:0;not null"`
	ShippingCost         float64          `gorm:"type:decimal(10,2);default:0;not null"`
	OrderItems           []OrderItem      `gorm:"foreignKey:OrderID"`
	OrderShippingAddress *ShippingAddress `gorm:"foreignKey:OrderID"`
	PaymentDetails       *Payment         `gorm:"-"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {

	o.OrderIdentifier = utils.GenerateOrderID()

	return nil

}
