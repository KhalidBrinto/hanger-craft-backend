package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

// CreateCategory creates a new category
func AddBrand(c *gin.Context) {

	var brand *models.Brand

	// Bind the incoming JSON to the Category struct
	if err := c.BindJSON(&brand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert the category into the database
	if err := config.DB.Create(&brand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the newly created category
	c.JSON(http.StatusCreated, gin.H{"message": "brand added successfully"})
}

// GetCategories retrieves all categories with their products
func GetBrands(c *gin.Context) {
	var brands []*models.Brand

	model := config.DB.Model(&brands).Preload("Products").Preload("Logo")

	pg := paginate.New()
	page := pg.With(model).Request(c.Request).Response(&brands)

	if page.Error {
		c.JSON(http.StatusInternalServerError, gin.H{"error": page.ErrorMessage})
		return
	}

	// Return the categories list
	c.JSON(http.StatusOK, page)
}

// UpdateCategory updates a category by its ID
func UpdateBrand(c *gin.Context) {

	brandID := c.Param("id")
	var brand *models.Brand

	// Fetch the category from the database
	if err := config.DB.First(&brand, brandID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Bind the updated data to the category
	if err := c.ShouldBindJSON(&brand); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the updated category
	if err := config.DB.Save(&brand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated category
	c.JSON(http.StatusOK, gin.H{"message": "Brand updated"})
}

// DeleteCategory deletes a category by its ID
func DeleteBrand(c *gin.Context) {

	brandID := c.Param("id")
	var brand *models.Category

	// Fetch the category from the database
	if err := config.DB.First(&brand, brandID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Delete the category from the database
	if err := config.DB.Delete(&brand).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success message
	c.JSON(http.StatusNoContent, gin.H{"message": "Brand deleted successfully"})
}
