package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateCategory creates a new category
func CreateCategory(c *gin.Context) {
	var category *models.Category

	// Bind the incoming JSON to the Category struct
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the category into the database
	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the newly created category
	c.JSON(http.StatusOK, category)
}

// GetCategories retrieves all categories with their products
func GetCategories(c *gin.Context) {
	var categories []*models.Category

	// Use Preload to load associated Products for each category
	if err := config.DB.Preload("Products").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the categories list
	c.JSON(http.StatusOK, categories)
}

// GetCategory retrieves a single category by its ID
func GetCategory(c *gin.Context) {
	categoryID := c.Param("id")
	var category *models.Category

	// Use Preload to load associated Products for this category
	if err := config.DB.Preload("Products").First(&category, categoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Return the category
	c.JSON(http.StatusOK, category)
}

// UpdateCategory updates a category by its ID
func UpdateCategory(c *gin.Context) {
	categoryID := c.Param("id")
	var category *models.Category

	// Fetch the category from the database
	if err := config.DB.First(&category, categoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Bind the updated data to the category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the updated category
	if err := config.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated category
	c.JSON(http.StatusOK, category)
}

// DeleteCategory deletes a category by its ID
func DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")
	var category *models.Category

	// Fetch the category from the database
	if err := config.DB.First(&category, categoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Delete the category from the database
	if err := config.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
