package controllers

import (
	"backend/config"
	"backend/models"
	"backend/serializers"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
