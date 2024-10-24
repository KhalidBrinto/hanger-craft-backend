package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	PaymentMethod  string  `gorm:"size:50;not null;check:payment_method IN ('card', 'bkash', 'rocket', 'nagad', 'cash_on_delivery')"`
	PaymentStatus  string  `gorm:"size:50;not null;check:payment_status IN ('pending', 'completed', 'failed')"`
	Amount         float64 `gorm:"not null"`
	TransanctionID *string
	PaymentDate    time.Time
	OrderID        uint  `gorm:"not null"`
	Order          Order `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}
