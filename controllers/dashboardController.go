package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetStats returns total orders, total revenue, and total customers
func GetStats(c *gin.Context) {
	var response struct {
		TotalOrder    int
		TotalRevenue  int
		TotalCustomer int
	}

	// Count total orders
	if err := config.DB.Model(&models.Order{}).
		Select("count(id) as total_order, sum(total_price) as total_revenue, (select count(users.id) from users where role = 'customer') as total_customer").
		Find(&response).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}

	// Return stats
	c.JSON(http.StatusOK, response)
}

// GetMonthlySales returns total sales for each month of the current year
func GetMonthlySales(c *gin.Context) {
	var monthlySales []struct {
		Month int     `json:"month"`
		Sales float64 `json:"sales"`
	}

	currentYear := time.Now().Year()

	// Query to get monthly sales for the current year
	if err := config.DB.Raw(`
		SELECT 
			EXTRACT(MONTH FROM created_at) AS month, 
			SUM(total_price) AS sales
		FROM orders
		WHERE EXTRACT(YEAR FROM created_at) = ?
		GROUP BY month
		ORDER BY month`, currentYear).Scan(&monthlySales).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve monthly sales"})
		return
	}

	// Return the result
	c.JSON(http.StatusOK, gin.H{"monthly_sales": monthlySales})
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
			TO_CHAR(DATE_TRUNC('month', created_at), 'Mon YYYY') AS month, 
			SUM(total_price) AS revenue
		FROM orders
		WHERE created_at BETWEEN ? AND ?
		GROUP BY month
		ORDER BY month ASC`, startDate, now).Scan(&yearlyRevenue).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve yearly revenue"})
		return
	}

	// Return the result
	c.JSON(http.StatusOK, gin.H{"yearly_revenue": yearlyRevenue})
}
