package controllers

import (
	"backend/config"
	"backend/models"
	"backend/serializers"
	"backend/utils"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

// CreateOrder creates a new order with order items and updates the inventory
func CreateOrder(c *gin.Context) {
	var order *models.Order
	var shipping_option *models.ShippingOptions

	// Bind JSON request to order struct
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order.UserID = c.GetUint("user_id")
	order.OrderStatus = "pending"
	order.ItemPrice = 0.0

	if err := config.DB.Where("payment_method = ?", order.PaymentDetails.PaymentMethod).Find(&shipping_option).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "invalid payment method"})
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

	fmt.Printf("lenth of order items: %d\n", len(order.OrderItems))

	// Loop through the order items and create them, also update inventory for each product
	for _, item := range order.OrderItems {
		order.ItemPrice += item.PriceAtPurchase * float64(item.Quantity)

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
		inventory.InOpen += item.Quantity
		inventory.ChangeType = "purchase"
		inventory.ChangeDate = time.Now()

		if err := tx.Save(&inventory).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
			return
		}

	}

	order.TotalPrice = order.ItemPrice - order.DiscountAmount + order.ShippingCost

	order.PaymentDetails.OrderID = order.ID
	order.PaymentDetails.Amount = order.TotalPrice
	order.PaymentDetails.TransanctionID = toPtr(utils.GenerateTransactionID())
	order.PaymentDetails.PaymentStatus = "pending"

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
		return
	}

	if err := tx.Save(&order.PaymentDetails).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment details"})
		return
	}

	// Commit the transaction
	tx.Commit()

	// Return the created order and inventory updates
	c.JSON(http.StatusOK, gin.H{"message": "order created successfully", "OrderID": order.OrderIdentifier})
}

// GetOrder retrieves an order by ID along with its items
func GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")
	var order *serializers.OrderResponse

	// Preload OrderItems to include them in the response
	if err := config.DB.Model(&models.Order{}).Preload("User").Preload("PaymentDetails").Preload("OrderItems.Product").Preload("OrderShippingAddress").First(&order, orderID).Error; err != nil {
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

func GetOrders(c *gin.Context) {
	var order []*serializers.OrderResponse
	var model *gorm.DB

	if c.GetString("role") == "admin" {

		// Preload OrderItems to include them in the response
		model = config.DB.Model(&models.Order{}).Preload("User").Preload("PaymentDetails").Preload("OrderItems.Product").Preload("OrderShippingAddress").Order("created_at DESC")

	} else {
		// Preload OrderItems to include them in the response
		model = config.DB.Model(&models.Order{}).Preload("User").Preload("OrderItems.Product").Preload("OrderShippingAddress").Where("user_id = ?", c.GetUint("user_id")).Order("created_at DESC")

	}

	pg := paginate.New()
	page := pg.With(model).Request(c.Request).Response(&order)

	if page.Error {
		c.JSON(http.StatusInternalServerError, gin.H{"error": page.ErrorMessage})
		return
	}

	// Return the order with its items
	c.JSON(http.StatusOK, &page)
}

// DispatchOrder updates an order status to shipped by its ID
func DispatchOrder(c *gin.Context) {

	orderID := c.Param("id")

	// Fetch the category from the database
	if err := config.DB.Model(&models.Order{}).Where("id = ?", orderID).Update("order_status", "shipped").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Return the updated category
	c.JSON(http.StatusOK, gin.H{"message": "order dispatched"})
}

// CancelOrder updates an order status to cancelled by its ID
func CancelOrder(c *gin.Context) {

	orderID := c.Param("id")

	// Fetch the category from the database
	if err := config.DB.Model(&models.Order{}).Where("id = ?", orderID).Update("order_status", "cancelled").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Return the updated category
	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})
}

// CancelOrder updates an order status to cancelled by its ID
func UpdateOrderStatus(c *gin.Context) {

	orderID := c.Param("id")

	var payload struct {
		OrderStatus string `binding:"required"`
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&models.Order{}).Where("id = ?", orderID).Update("order_status", payload.OrderStatus).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Return the updated category
	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})
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

// RestockProduct adds stock for a given product
func GetInventory(c *gin.Context) {
	var inventory []*serializers.InventoryResponse

	// Preload OrderItems to include them in the response
	if err := config.DB.Model(&models.Inventory{}).Preload("Product").Select("inventories.*, (stock_level-in_open) as available_quantity").Find(&inventory).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No inventory found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, inventory)
}
