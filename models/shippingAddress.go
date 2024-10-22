package models

import "gorm.io/gorm"

type ShippingAddress struct {
	gorm.Model
	UserID       uint   `gorm:"not null"`
	AddressLine1 string `gorm:"size:255;not null"`
	AddressLine2 string `gorm:"size:255"`
	City         string `gorm:"size:100;not null"`
	State        string `gorm:"size:100"`
	PostalCode   string `gorm:"size:20;not null"`
	Country      string `gorm:"size:100;not null"`

	User User
}
