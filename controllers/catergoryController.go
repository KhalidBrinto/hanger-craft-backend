package controllers

import (
	"backend/config"
	"backend/models"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

// CreateCategory creates a new category
func CreateCategory(c *gin.Context) {

	var category *models.Category

	// Bind the incoming JSON to the Category struct
	if err := c.BindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CategoryType must be in ['parent', 'child']"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the newly created category
	c.JSON(http.StatusCreated, gin.H{"message": "category added successfully"})
	// Insert the category into the database
	// if err := config.DB.Create(&category).Error; err != nil {
	// 	if errors.Is(err, gorm.ErrCheckConstraintViolated) {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "CategoryType must be in ['parent', 'child']"})
	// 		return
	// 	}
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
}

// GetCategories retrieves all categories with their products
func GetCategories(c *gin.Context) {
	var categories []*models.Category
	querystring := ""
	if c.Query("type") != "" {
		querystring = "category_type = '" + c.Query("type") + "'"
	}

	// Use Preload to load associated Products for each category
	if err := config.DB.Preload("Products").Preload("Image").Where(querystring).Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the categories list
	c.JSON(http.StatusOK, categories)
}

func GetNestedCategories(c *gin.Context) {

	type Category struct {
		gorm.Model
		Name         null.String `gorm:"size:100;not null"`
		CategoryType null.String `gorm:"size:100;not null;check:category_type IN ('parent', 'child')"`
		ParentID     *uint
		// Products     []Product `gorm:"foreignKey:CategoryID"`
		SubCategory json.RawMessage
	}
	var categories []*Category

	model := config.DB.Model(&categories).
		Select(`categories.*, 
				COALESCE(
					json_agg(
						json_build_object(
						'ID', subcategories.id,
						'Name', subcategories.name,
						'ParentID', subcategories.parent_id
						)
					)FILTER (WHERE subcategories.id IS NOT NULL),
            		'[]'
				) AS sub_category
			`).
		Joins("LEFT JOIN categories AS subcategories ON subcategories.parent_id = categories.id").
		Where("categories.parent_id is null").
		Group("categories.id").
		Find(&categories)

	if model.Error != nil {
		if errors.Is(model.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "No categories found"})
			return

		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": model.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, &categories)
}

func GetSubCategories(c *gin.Context) {
	parentID := c.Param("parent_id")
	var categories []*models.Category

	// Use Preload to load associated Products for each category
	if err := config.DB.Preload("Products").Where("parent_id = ?", parentID).Find(&categories).Error; err != nil {
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
