package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	PaymentMethod  string  `gorm:"size:50;not null;check:payment_method IN ('card', 'bkash', 'rocket', 'nagad', 'cash_on_delivery')"`
	PaymentStatus  string  `gorm:"size:50;not null;check:payment_status IN ('pending', 'completed', 'failed')"`
	Amount         float64 `gorm:"type:decimal(10,2);not null"`
	TransanctionID *string `gorm:"size:11;not null"`
	PaymentDate    *time.Time
	OrderID        uint  `gorm:"not null"`
	Order          Order `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

type PaymentOption struct {
	gorm.Model
	MerchantID    *string `gorm:"size:255"`
	PaymentMethod string  `gorm:"size:50;not null;check:payment_method IN ('card', 'bkash', 'rocket', 'nagad', 'cash_on_delivery')"`
	Status        bool    `gorm:"not null;default:false"`
	APIKey        *string `gorm:"type:text"` // API key for authenticating requests
	APISecret     *string `gorm:"type:text"`
}
