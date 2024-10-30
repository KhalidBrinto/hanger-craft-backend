package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
)

// CreateProduct creates a new product
func CreateProduct(c *gin.Context) {
	var product *models.Product

	if err := c.BindJSON(&product); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&product).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&models.Inventory{
		ProductID:  product.ID,
		StockLevel: int(product.Stock),
		InOpen:     0,
		ChangeType: "restock",
		ChangeDate: time.Now(),
	}).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product added successfully"})
}

// GetProducts retrieves all products with their category and reviews
func GetProducts(c *gin.Context) {
	var params utils.Parameters
	if c.Bind(&params) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind identifier parameters."})
		return
	}
	querystring := utils.ProductQueryParameterToMap(params)

	type Inventory struct {
		ProductID  uint           `gorm:"not null" json:"-"`
		Product    models.Product `gorm:"foreignKey:ProductID" json:"-"`
		StockLevel int            `gorm:"not null"`
	}
	type Product struct {
		gorm.Model
		Name         string          `gorm:"size:150;not null"`
		Description  string          `gorm:"type:text"`
		SKU          string          `gorm:"size:150;not null;unique;index"`
		Barcode      *string         `gorm:"size:150"`
		Price        float64         `gorm:"type:decimal(10,2);not null"`
		Currency     string          `gorm:"size:3; not null"`
		Images       pq.StringArray  `gorm:"type:varchar[]"`
		CategoryID   uint            `gorm:"not null"`
		Category     models.Category `gorm:"foreignKey:CategoryID"`
		Status       *string         `gorm:"not null;check:status IN ('published', 'unpublished')"`
		Inventory    *Inventory      `gorm:"foreignKey:ProductID"`
		TotalReviews int
		Rating       int
	}

	var products []*Product
	var model *gorm.DB

	model = config.DB.Model(&products).Preload("Category").Preload("Inventory").
		Select(`products.*, 
				count(reviews.id) as total_reviews,
				AVG(reviews.rating)::int as rating
			`).
		Joins("LEFT JOIN reviews ON products.id = reviews.product_id").
		Where(querystring).
		Group("products.id")

	pg := paginate.New()
	page := pg.With(model).Request(c.Request).Response(&products)

	if page.Error {
		c.JSON(http.StatusInternalServerError, gin.H{"error": page.ErrorMessage})
		return
	}

	c.JSON(http.StatusOK, &page)
}
func GetNewArrivalProducts(c *gin.Context) {
	var params utils.Parameters
	if c.Bind(&params) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind identifier parameters."})
		return
	}
	querstring := utils.ProductQueryParameterToMap(params)
	type Inventory struct {
		ProductID  uint           `gorm:"not null" json:"-"`
		Product    models.Product `gorm:"foreignKey:ProductID" json:"-"`
		StockLevel int            `gorm:"not null"`
	}
	type Product struct {
		gorm.Model
		Name         string          `gorm:"size:150;not null"`
		Description  string          `gorm:"type:text"`
		SKU          string          `gorm:"size:150;not null;unique;index"`
		Barcode      *string         `gorm:"size:150"`
		Price        float64         `gorm:"type:decimal(10,2);not null"`
		Currency     string          `gorm:"size:3; not null"`
		Images       pq.StringArray  `gorm:"type:varchar[]"`
		CategoryID   uint            `gorm:"not null"`
		Category     models.Category `gorm:"foreignKey:CategoryID"`
		Status       *string         `gorm:"not null;check:status IN ('published', 'unpublished')"`
		Inventory    *Inventory      `gorm:"foreignKey:ProductID"`
		TotalReviews int
		Rating       int
	}

	var products []*Product
	var model *gorm.DB

	model = config.DB.Model(&products).Preload("Category").Preload("Inventory").
		Select(`products.*, 
				count(reviews.id) as total_reviews,
				AVG(reviews.rating)::int as rating
			`).
		Joins("LEFT JOIN reviews ON products.id = reviews.product_id").
		Where(querstring).
		Group("products.id").
		Order("products.created_at DESC")

	pg := paginate.New()
	page := pg.With(model).Request(c.Request).Response(&products)

	if page.Error {
		c.JSON(http.StatusInternalServerError, gin.H{"error": page.ErrorMessage})
		return
	}

	c.JSON(http.StatusOK, &page)
}
func GetTrendingProducts(c *gin.Context) {
	var params utils.Parameters
	if c.Bind(&params) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind identifier parameters."})
		return
	}
	querstring := utils.ProductQueryParameterToMap(params)
	type Inventory struct {
		ProductID  uint           `gorm:"not null" json:"-"`
		Product    models.Product `gorm:"foreignKey:ProductID" json:"-"`
		StockLevel int            `gorm:"not null"`
	}
	type Product struct {
		gorm.Model
		Name         string          `gorm:"size:150;not null"`
		Description  string          `gorm:"type:text"`
		SKU          string          `gorm:"size:150;not null;unique;index"`
		Barcode      *string         `gorm:"size:150"`
		Price        float64         `gorm:"type:decimal(10,2);not null"`
		Currency     string          `gorm:"size:3; not null"`
		Images       pq.StringArray  `gorm:"type:varchar[]"`
		CategoryID   uint            `gorm:"not null"`
		Category     models.Category `gorm:"foreignKey:CategoryID"`
		Status       *string         `gorm:"not null;check:status IN ('published', 'unpublished')"`
		Inventory    *Inventory      `gorm:"foreignKey:ProductID"`
		TotalReviews int
		Rating       int
	}

	var products []*Product
	var model *gorm.DB

	model = config.DB.Model(&products).Preload("Category").Preload("Inventory").
		Select(`products.*, 
				count(reviews.id) as total_reviews,
				AVG(reviews.rating)::int as rating
			`).
		Joins("LEFT JOIN reviews ON products.id = reviews.product_id").
		Joins("LEFT JOIN order_items on products.id = order_items.product_id").
		Where(querstring).
		Group("products.id").
		Order("COUNT(distinct order_items.order_id) DESC")

	pg := paginate.New()
	page := pg.With(model).Request(c.Request).Response(&products)

	if page.Error {
		c.JSON(http.StatusInternalServerError, gin.H{"error": page.ErrorMessage})
		return
	}

	c.JSON(http.StatusOK, &page)
}

// GetProduct retrieves a single product by its ID
func GetProduct(c *gin.Context) {
	productID := c.Param("id")
	var product *models.Product

	if err := config.DB.Preload("Category").First(&product, productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct updates a product by its ID
func UpdateProduct(c *gin.Context) {
	productID := c.Param("id")
	var product *models.Product

	if err := config.DB.First(&product, productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct deletes a product by its ID
func DeleteProduct(c *gin.Context) {
	productID := c.Param("id")
	var product *models.Product

	if err := config.DB.First(&product, productID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// CreateProductAttribute creates a new product attribute
func CreateProductAttribute(c *gin.Context) {
	var productAttribute *models.ProductAttribute

	if err := c.ShouldBindJSON(&productAttribute); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&productAttribute).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productAttribute)
}

// GetProductAttributes retrieves all product attributes
func GetProductAttributes(c *gin.Context) {
	var productAttributes []*models.ProductAttribute

	if err := config.DB.Preload("Product").Find(&productAttributes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productAttributes)
}

// UpdateProductAttribute updates a product attribute by its ID
func UpdateProductAttribute(c *gin.Context) {
	attributeID := c.Param("id")
	var productAttribute *models.ProductAttribute

	if err := config.DB.First(&productAttribute, attributeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product attribute not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := c.ShouldBindJSON(&productAttribute); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&productAttribute).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productAttribute)
}

// DeleteProductAttribute deletes a product attribute by its ID
func DeleteProductAttribute(c *gin.Context) {
	attributeID := c.Param("id")
	var productAttribute *models.ProductAttribute

	if err := config.DB.First(&productAttribute, attributeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product attribute not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := config.DB.Delete(&productAttribute).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product attribute deleted successfully"})
}
