package models

import (
	"time"

	"gorm.io/gorm"
)

type Coupon struct {
	gorm.Model
	Code              string     `gorm:"size:50;unique;not null"`                                         // Unique coupon code
	Description       string     `gorm:"type:text"`                                                       // Description of the coupon
	DiscountType      string     `gorm:"size:20;not null;check:discount_type IN ('percentage', 'fixed')"` // Type of discount: 'percentage' or 'fixed'
	DiscountValue     float64    `gorm:"type:numeric(10,2);not null"`                                     // Discount value (percentage or fixed amount)
	MinOrderValue     *float64   `gorm:"type:numeric(10,2)"`                                              // Minimum order value required to use the coupon
	MaxDiscountValue  *float64   `gorm:"type:numeric(10,2)"`                                              // Max discount for percentage-based coupons
	UsageLimit        *int       // Total times this coupon can be used
	UsageLimitPerUser int        `gorm:"default:1"` // Times each user can use the coupon
	StartDate         time.Time  `gorm:"not null"`  // Start date for coupon validity
	ExpirationDate    *time.Time // Expiration date for coupon validity
	IsActive          bool       `gorm:"default:true"` // Whether the coupon is active
}

type CouponUsageHistory struct {
	ID        uint      `gorm:"primaryKey"`
	CouponID  uint      `gorm:"not null"` // Reference to Coupon
	Category  Coupon    `gorm:"foreignKey:CouponID"`
	UserID    uint      `gorm:"not null"` // Reference to the user who used the coupon
	User      User      `gorm:"foreignKey:UserID"`
	UsedAt    time.Time `gorm:"autoCreateTime"` // Timestamp of when the coupon was used
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
