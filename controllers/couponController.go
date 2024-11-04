package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Get all coupons
func GetCoupons(c *gin.Context) {
	var coupons []models.Coupon
	if err := config.DB.Find(&coupons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch coupons"})
		return
	}
	c.JSON(http.StatusOK, coupons)
}

// Get a coupon by ID
func GetCoupon(c *gin.Context) {
	id := c.Param("id")
	var coupon models.Coupon
	if err := config.DB.First(&coupon, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}
	c.JSON(http.StatusOK, coupon)
}

// Create a new coupon
func CreateCoupon(c *gin.Context) {
	var coupon models.Coupon
	if err := c.ShouldBindJSON(&coupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create coupon"})
		return
	}
	c.JSON(http.StatusCreated, coupon)
}

// Update a coupon
func UpdateCoupon(c *gin.Context) {
	id := c.Param("id")
	var coupon models.Coupon
	if err := config.DB.First(&coupon, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}
	if err := c.ShouldBindJSON(&coupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Save(&coupon).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Coupon updated"})
}

// Delete a coupon
func DeleteCoupon(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Coupon{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coupon"})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Coupon deleted"})
}

// Apply a coupon
func ApplyCoupon(c *gin.Context, CouponCode string, userID uint) *models.Coupon {

	var coupon models.Coupon
	if err := config.DB.Where("code = ? AND is_active = true", CouponCode).First(&coupon).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found or inactive"})
		return nil
	}

	// Check if coupon has expired
	if coupon.ExpirationDate != nil && time.Now().After(*coupon.ExpirationDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coupon has expired"})
		return nil
	}

	// Check usage limits
	var usageCount int64
	config.DB.Model(&models.CouponUsageHistory{}).Where("coupon_id = ?", coupon.ID).Count(&usageCount)
	if coupon.UsageLimit != nil && usageCount >= int64(*coupon.UsageLimit) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Coupon usage limit reached"})
		return nil
	}

	var userUsageCount int64
	config.DB.Model(&models.CouponUsageHistory{}).Where("coupon_id = ? AND user_id = ?", coupon.ID, userID).Count(&userUsageCount)
	if userUsageCount >= int64(coupon.UsageLimitPerUser) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User coupon usage limit reached"})
		return nil
	}

	// Log coupon usage
	usageHistory := models.CouponUsageHistory{
		CouponID: coupon.ID,
		UserID:   userID,
		UsedAt:   time.Now(),
	}
	if err := config.DB.Create(&usageHistory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log coupon usage"})
		return nil
	}

	return &coupon
}
