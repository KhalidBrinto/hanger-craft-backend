package models

import (
	"time"

	"gorm.io/gorm"
)

type Inventory struct {
	gorm.Model
	ProductID  uint    `gorm:"not null"`
	Product    Product `gorm:"foreignKey:ProductID"`
	StockLevel int     `gorm:"not null"`
	InOpen     int     `gorm:"not null"`
	ChangeType string  `gorm:"size:50;not null;check:change_type IN ('restock', 'purchase')"`
	ChangeDate time.Time
}
