package controllers

import (
	"backend/config"
	"backend/models"
	"errors"
	"fmt"
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

	fmt.Println(order.ID)

	// Loop through the order items and create them, also update inventory for each product
	for _, item := range order.OrderItems {
		// Create order items
		if err := tx.Create(&models.OrderItem{
			OrderID:         order.ID,
			ProductID:       item.ProductID,
			Quantity:        item.Quantity,
			PriceAtPurchase: item.PriceAtPurchase,
		}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order items"})
			return
		}

		// Fetch the existing inventory record for the product
		var inventory models.Inventory
		if err := tx.Where("product_id = ?", item.ProductID).First(&inventory).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inventory"})
			return
		}

		// Check if there's enough stock to fulfill the order
		if inventory.StockLevel < item.Quantity+inventory.InOpen {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough stock available"})
			return
		}

		// Update inventory: subtract ordered quantity from InOpen and StockLevel
		// inventory.StockLevel -= item.Quantity
		inventory.InOpen += item.Quantity

		// Log inventory change
		inventoryUpdate := models.Inventory{
			ProductID:  item.ProductID,
			InOpen:     inventory.InOpen, // Updated in-open quantity
			ChangeType: "purchase",       // Record the change type
			ChangeDate: time.Now(),
		}

		if err := tx.Save(&inventory).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
			return
		}

		if err := tx.Create(&inventoryUpdate).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log inventory change"})
			return
		}

		inventoryUpdates = append(inventoryUpdates, &inventoryUpdate)
	}

	// Commit the transaction
	tx.Commit()

	// Return the created order and inventory updates
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
	var inventoryRequest *models.Inventory

	// Bind JSON request to inventoryRequest struct
	if err := c.ShouldBindJSON(&inventoryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingInventory models.Inventory

	// Check if an inventory record already exists for the given product
	if err := config.DB.Where("product_id = ?", inventoryRequest.ProductID).First(&existingInventory).Error; err != nil {
		// If no existing record, create a new one
		if errors.Is(err, gorm.ErrRecordNotFound) {
			inventoryRequest.ChangeType = "restock"
			inventoryRequest.ChangeDate = time.Now()

			// Create new inventory record
			if err := config.DB.Create(&inventoryRequest).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create inventory record"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Stock added successfully", "inventory": inventoryRequest})
		} else {
			// Handle other database errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// If the inventory record exists, update stock levels
		existingInventory.StockLevel += inventoryRequest.StockLevel
		existingInventory.ChangeType = "restock"
		existingInventory.ChangeDate = time.Now() // Add new stock to the current stock level

		// Save the updated stock levels
		if err := config.DB.Save(&existingInventory).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory record"})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully", "inventory": existingInventory})
	}
}
