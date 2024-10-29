package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morkid/paginate"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Helper function to create a pointer from a string value
func toPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func RegisterCustomer(c *gin.Context) {
	var input struct {
		Name        string  `json:"name" binding:"required"`
		Email       string  `json:"email" binding:"required,email"`
		Address     *string `json:"address"`
		Password    string  `json:"password" binding:"required"`
		PhoneNumber string  `json:"phone_number"`
	}

	// Bind the JSON input to the struct
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: toPtr(string(passwordHash)),
		PhoneNumber:  toPtr(input.PhoneNumber),
		Address:      input.Address,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func UpdateUser(c *gin.Context) {
	user_id := c.GetString("user_id")
	var user *models.User

	// Find the payment by ID
	if err := config.DB.First(&user, user_id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// LoginUser handles user login
func LoginUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Bind the JSON input
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user by email
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user"})
		return
	}

	// Compare hashed password with input password
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "access_token": token})
}

func GetCustomers(c *gin.Context) {
	var customers []struct {
		gorm.Model
		Name         string  `gorm:"size:100;not null"`
		Email        string  `gorm:"size:100;unique;not null"`
		PhoneNumber  *string `gorm:"size:15"`
		Address      json.RawMessage
		Role         string `gorm:"size:20;default:'customer';not null"`
		TotalOrders  int
		LastPurchase *time.Time
	}

	// Use Preload to load associated Products for each category
	model := config.DB.Model(&models.User{}).
		Select(`
		users.*,
		count(orders.id) as total_orders,
		max(orders.created_at) as last_purchase
	`).
		Joins("LEFT JOIN orders ON users.id = orders.user_id").
		Where("role = ?", "customer").
		Group("users.id").
		Order("total_orders DESC").
		Find(&customers)

	pg := paginate.New()
	page := pg.With(model).Request(c.Request).Response(&customers)

	if page.Error {
		c.JSON(http.StatusInternalServerError, gin.H{"error": page.ErrorMessage})
		return
	}

	// Return the categories list
	c.JSON(http.StatusOK, &page)
}

func DeleteCustomer(c *gin.Context) {
	userID := c.GetUint("user_id")
	var customer *models.User

	if err := config.DB.First(&customer, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	if err := config.DB.Delete(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Return success message
	c.JSON(http.StatusNoContent, gin.H{"message": "User deleted successfully"})
}

func DeleteUserByID(c *gin.Context) {
	userID := c.Param("id")
	var user *models.User

	if err := config.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		}
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Return success message
	c.JSON(http.StatusNoContent, gin.H{"message": "User deleted successfully"})
}
