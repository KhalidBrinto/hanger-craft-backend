package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateOrder creates a new order with order items and updates the inventory
func CreateOrder(c *gin.Context) {
	var order *models.Order
	var inventoryUpdates []*models.Inventory

	// Bind JSON request to order struct
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a database transaction
	tx := config.DB.Begin()

	// Insert the order in the database
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Loop through the order items and create them, also update inventory for each product
	for _, item := range order.OrderItems {
		// Create order items
		item.OrderID = order.ID // Assign order ID
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order items"})
			return
		}

		// Update inventory by subtracting the ordered quantity
		inventoryUpdate := models.Inventory{
			ProductID:      item.ProductID,
			QuantityChange: -item.Quantity, // Deduct stock
			ChangeType:     "purchase",
			ChangeDate:     time.Now(),
		}
		if err := tx.Create(&inventoryUpdate).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
			return
		}

		inventoryUpdates = append(inventoryUpdates, &inventoryUpdate)
	}

	// Commit the transaction
	tx.Commit()

	// Return the created order
	c.JSON(http.StatusOK, gin.H{"order": order, "inventory_updates": inventoryUpdates})
}

// GetOrder retrieves an order by ID along with its items
func GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	var order *models.Order

	// Preload OrderItems to include them in the response
	if err := config.DB.Preload("OrderItems.Product").First(&order, orderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Return the order with its items
	c.JSON(http.StatusOK, order)
}

// RestockProduct adds stock for a given product
func RestockProduct(c *gin.Context) {
	var inventory *models.Inventory

	// Bind JSON request to inventory struct
	if err := c.ShouldBindJSON(&inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the change type is 'restock'
	inventory.ChangeType = "restock"
	inventory.ChangeDate = time.Now()

	// Create an inventory change record
	if err := config.DB.Create(&inventory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "Stock added successfully", "inventory": inventory})
}
