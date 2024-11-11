package controllers

import (
	"backend/config"
	"backend/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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
		ID       uint        `json:"id"`
		Name     string      `json:"name"`
		ParentID *uint       `json:"parent_id"` // Nullable for top-level categories
		Level    int         `json:"level"`
		Path     []int       `json:"-"` // Path to use for ordering; ignore in JSON
		Children []*Category `json:"children,omitempty"`
	}
	var categories []*Category

	query := `
		WITH RECURSIVE category_hierarchy AS (
			SELECT id, name, parent_id, 1 AS level, ARRAY[id] AS path
			FROM categories
			WHERE parent_id IS NULL
			UNION ALL
			SELECT c.id, c.name, c.parent_id, ch.level + 1, ch.path || c.id
			FROM categories c
			JOIN category_hierarchy ch ON c.parent_id = ch.id
		)
		SELECT id, name, parent_id, level, path
		FROM category_hierarchy
		ORDER BY path;`

	if err := config.DB.Raw(query).Scan(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	buildCategoryTree := func(categories []*Category) []*Category {
		categoryMap := make(map[uint]*Category)
		var rootCategories []*Category

		for i := range categories {
			category := categories[i]           // Current category
			categoryMap[category.ID] = category // Store in map

			if category.ParentID == nil {
				rootCategories = append(rootCategories, category) // Add to root if no parent
			} else {
				parent := categoryMap[*category.ParentID]
				parent.Children = append(parent.Children, category) // Append to parent’s children
			}
		}

		return rootCategories
	}
	nestedCategories := buildCategoryTree(categories)
	c.JSON(http.StatusOK, nestedCategories)
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
