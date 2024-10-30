package controllers

import (
	"backend/config"
	"backend/models"
	"backend/serializers"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

// CreateReview handles creating a new review for a product
func CreateReview(c *gin.Context) {
	var review *models.Review

	// Bind the JSON request to the Review struct
	if err := c.BindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review.UserID = c.GetUint("user_id")

	// Set the review creation date if it's not provided
	if review.CreatedAt.IsZero() {
		review.CreatedAt = time.Now()
	}

	// Create the review in the database
	if err := config.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create review"})
		return
	}

	// Return the created review
	c.JSON(http.StatusOK, gin.H{"message": "review recorded successfully"})
}

// GetReview retrieves a review by its ID
func GetReview(c *gin.Context) {
	reviewID := c.Param("id")
	var review *models.Review

	// Find the review by ID and preload the associated user and product
	if err := config.DB.Preload("User").Preload("Product").First(&review, reviewID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Return the review with the associated user and product
	c.JSON(http.StatusOK, review)
}

func GetCustomerReview(c *gin.Context) {
	type Product struct {
		ID          uint            `gorm:"primarykey"`
		Name        string          `gorm:"size:150;not null"`
		Description string          `gorm:"type:text"`
		SKU         string          `gorm:"size:150;not null;unique;index"`
		Barcode     *string         `gorm:"size:150"`
		Price       float64         `gorm:"type:decimal(10,2);not null"`
		Currency    string          `gorm:"size:3; not null"`
		Images      pq.StringArray  `gorm:"type:varchar[]"`
		CategoryID  uint            `gorm:"not null"`
		Category    models.Category `gorm:"foreignKey:CategoryID"`
	}
	type User struct {
		ID          uint            `gorm:"primarykey"`
		Name        string          `gorm:"size:100;not null"`
		Email       string          `gorm:"size:100;unique;not null"`
		Address     json.RawMessage `gorm:"type:jsonb"`
		PhoneNumber *string         `gorm:"size:15"`
	}
	var reviews []*struct {
		gorm.Model
		UserID    uint    `gorm:"not null" json:"-"`
		User      User    `gorm:"foreignKey:UserID"`
		ProductID uint    `gorm:"not null" json:"-"`
		Product   Product `gorm:"foreignKey:ProductID"`
		Rating    int     `gorm:"check:rating >= 1 AND rating <= 5"`
		Comment   string  `gorm:"type:text"`
	}

	// Find the review by ID and preload the associated user and product
	model := config.DB.Model(&models.Review{}).Preload("User").Preload("Product").Preload("Product.Category").Order("created_at DESC")

	pg := paginate.New()
	page := pg.With(model).Request(c.Request).Response(&reviews)

	if page.Error {
		c.JSON(http.StatusInternalServerError, gin.H{"error": page.ErrorMessage})
		return
	}

	// Return the review with the associated user and product
	c.JSON(http.StatusOK, &page)
}

// GetReviewsByProduct retrieves all reviews for a specific product
func GetReviewsByProduct(c *gin.Context) {
	productID := c.Param("product_id")
	var reviews []*serializers.ReviewResponse

	// Find all reviews for the specified product ID
	if err := config.DB.Model(&models.Review{}).Where("product_id = ?", productID).Preload("User").Order("created_at desc").Find(&reviews).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No reviews found for this product"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Return the list of reviews with associated users
	c.JSON(http.StatusOK, reviews)
}

// UpdateReview updates the rating or comment of a specific review
func UpdateReview(c *gin.Context) {
	reviewID := c.Param("id")
	var review *models.Review

	// Find the review by ID
	if err := config.DB.First(&review, reviewID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Bind the JSON request to the Review struct (for updating)
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the review in the database
	if err := config.DB.Save(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update review"})
		return
	}

	// Return the updated review
	c.JSON(http.StatusOK, gin.H{"review": review})
}

// DeleteReview deletes a review by its ID
func DeleteReview(c *gin.Context) {
	reviewID := c.Param("id")
	var review *models.Review

	// Find the review by ID
	if err := config.DB.First(&review, reviewID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Delete the review
	if err := config.DB.Delete(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete review"})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
