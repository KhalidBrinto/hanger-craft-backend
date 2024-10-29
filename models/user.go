package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string  `gorm:"size:100;not null"`
	Email        string  `gorm:"size:100;unique;not null"`
	Address      *string `gorm:"type:text"`
	PasswordHash *string `gorm:"size:255;not null" json:"password"`
	PhoneNumber  *string `gorm:"size:15"`
	Role         string  `gorm:"size:20;default:'customer';not null"`
}
