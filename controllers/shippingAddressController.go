package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateShippingAddress creates a shipping address for a specific order
func CreateShippingAddress(c *gin.Context) {
	var shippingAddress *models.ShippingAddress

	// Bind the JSON request to the ShippingAddress struct
	if err := c.ShouldBindJSON(&shippingAddress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the shipping address for the specific order
	if err := config.DB.Create(&shippingAddress).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shipping address"})
		return
	}

	// Return the created shipping address
	c.JSON(http.StatusOK, gin.H{"message": "shipping address created successfully"})
}
