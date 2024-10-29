package models

import (
	"encoding/json"

	"gorm.io/gorm"
)

type ShippingAddress struct {
	gorm.Model
	UserID       *uint
	OrderID      uint   `gorm:"not null"`
	AddressLine1 string `gorm:"size:255;not null"`
	AddressLine2 string `gorm:"size:255"`
	City         string `gorm:"size:100;not null"`
	State        string `gorm:"size:100"`
	PostalCode   string `gorm:"size:20;not null"`
	Country      string `gorm:"size:100;not null"`

	User  User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Order Order `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

type ShippingOptions struct {
	gorm.Model
	ShipperID               *string
	ShippingCarrier         *string         `gorm:"not null"`
	ShipFromAddress         json.RawMessage `gorm:"type:jsonb"`
	ShippingCost            float64         `gorm:"type:decimal(10,2)"`
	EstimatedDeliveryDayMin int             `gorm:"type:int"` // Minimum estimated days for delivery
	EstimatedDeliveryDayMax int             `gorm:"type:int"`
	PaymentMethod           *string         `gorm:"size:50;not null;check:payment_method IN ('card', 'bkash', 'rocket', 'nagad', 'cash_on_delivery')"`
}
