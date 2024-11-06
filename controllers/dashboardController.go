package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetStats returns total orders, total revenue, and total customers
func GetStats(c *gin.Context) {
	var response struct {
		TotalOrder     int
		TotalProducts  int
		TotalCustomers int
		TotalReviews   int
	}

	// Count total orders
	if err := config.DB.Model(&models.Order{}).
		Select(`
			count(id) as total_order, 
			(select count(products.id) from products where is_child = false AND deleted_at IS NULL) as total_products,
			(select count(users.id) from users where role = 'customer') as total_customers,
			(select count(reviews.id) from reviews where deleted_at IS NULL) as total_reviews
		`).
		Find(&response).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}

	// Return stats
	c.JSON(http.StatusOK, response)
}

func GetTopSellingProducts(c *gin.Context) {
	var response []struct {
		ProductID   uint
		ProductName string
		Quantity    int
	}

	// Count total orders
	if err := config.DB.Model(&models.OrderItem{}).
		Select(`
			product_id, 
			products.name as product_name, 
			SUM(quantity) as quantity
		`).
		Joins("LEFT JOIN products on products.id = product_id").
		Group("product_id, products.name").
		Order("quantity DESC").
		Find(&response).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top selling products"})
		return
	}

	// Return stats
	c.JSON(http.StatusOK, response)
}

// GetMonthlySales returns total sales for each month of the current year
func GetMonthlySales(c *gin.Context) {
	var monthlySales struct {
		Revenue    float64
		Total      int
		Completed  int
		Pending    int
		Cancelled  int
		LowStock   int
		OutOfStock int
	}

	currentMonth := int(time.Now().Month())
	currentYear := time.Now().Year()
	if c.Query("month") != "" {
		currentMonth, _ = strconv.Atoi(c.Query("month"))
	}
	if c.Query("year") != "" {
		currentYear, _ = strconv.Atoi(c.Query("year"))
	}

	// Query to get monthly sales for the current year
	if err := config.DB.Raw(`
		SELECT
			COUNT(id) as total,
			SUM(CASE WHEN order_status = 'delivered' THEN 1 ELSE 0 END) as completed,
			SUM(CASE WHEN order_status = 'pending' THEN 1 ELSE 0 END) as pending,
			SUM(CASE WHEN order_status = 'cancelled' THEN 1 ELSE 0 END) as cancelled
		FROM orders
		WHERE EXTRACT(MONTH FROM orders.created_at) = ? AND
		EXTRACT(YEAR FROM orders.created_at) = ?`, currentMonth, currentYear).Find(&monthlySales).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve monthly sales"})
		return
	}
	if err := config.DB.Raw(`
		SELECT
			SUM(amount) as revenue
		FROM payments
		WHERE payment_status = 'completed' AND
		EXTRACT(MONTH FROM payment_date) = ? AND
		EXTRACT(YEAR FROM payment_date) = ?`, currentMonth, currentYear).Find(&monthlySales).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve monthly sales"})
		return
	}
	if err := config.DB.Raw(`
		SELECT 
			SUM(CASE WHEN stock_level - in_open <= 10 THEN 1 ELSE 0 END) as low_stock,
			SUM(CASE WHEN stock_level - in_open = 0 THEN 1 ELSE 0 END) as out_of_stock
		FROM inventories`).
		Find(&monthlySales).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve inventory report"})
		return
	}

	if monthlySales.Total == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No data found"})
		return
	}
	// Return the result
	c.JSON(http.StatusOK, gin.H{
		"Revenue":    monthlySales.Revenue,
		"Completed":  monthlySales.Completed,
		"Pending":    monthlySales.Pending,
		"Cancelled":  monthlySales.Cancelled,
		"LowStock":   monthlySales.LowStock,
		"OutOfStock": monthlySales.OutOfStock,
	})
}

// GetYearlyRevenue returns the revenue for the past 12 months
func GetYearlyRevenue(c *gin.Context) {
	var yearlyRevenue []struct {
		Month   string  `json:"month"`
		Revenue float64 `json:"revenue"`
	}

	// Get the current month and the month one year ago
	now := time.Now()
	startDate := now.AddDate(-1, 0, 0)

	// Query to get the revenue for the last 12 months
	if err := config.DB.Raw(`
		SELECT 
			TO_CHAR(DATE_TRUNC('month', orders.created_at), 'Mon YYYY') AS month, 
			SUM(total_price) AS revenue
		FROM orders
		LEFT JOIN payments ON payments.order_id = orders.id
		WHERE orders.created_at BETWEEN ? AND ? AND payments.payment_status = 'completed'
		GROUP BY month
		ORDER BY month ASC`, startDate, now).Scan(&yearlyRevenue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve yearly revenue"})
		return
	}

	// Return the result
	c.JSON(http.StatusOK, gin.H{"yearly_revenue": yearlyRevenue})
}
