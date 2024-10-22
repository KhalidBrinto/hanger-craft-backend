package models

import "time"

type Review struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"not null"`
	User      User    `gorm:"foreignKey:UserID"`
	ProductID uint    `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
	Rating    int     `gorm:"check:rating >= 1 AND rating <= 5"`
	Comment   string  `gorm:"type:text"`
	CreatedAt time.Time
}
