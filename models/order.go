package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID      uint        `gorm:"not null"`
	User        User        `gorm:"foreignKey:UserID"`
	OrderStatus string      `gorm:"size:50;not null;check:order_status IN ('pending', 'shipped', 'delivered', 'cancelled')"`
	TotalPrice  float64     `gorm:"not null"`
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID"`
}
