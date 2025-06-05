package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/morkid/paginate"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

// CreateProduct creates a new product
func CreateProduct(c *gin.Context) {
	type Variation struct {
		Size     string
		Quantity int
	}
	type Attributes struct {
		Image     string
		Color     string
		Variation []Variation
	}
	var payload struct {
		Name        string  `gorm:"size:150;not null"`
		Description string  `gorm:"type:text"`
		SKU         string  `gorm:"size:150;not null;unique;index"`
		Barcode     *string `gorm:"size:150"`
		Price       float64 `gorm:"type:decimal(10,2);not null"`
		Currency    string  `gorm:"size:3; not null"`
		CategoryID  uint    `gorm:"not null"`
		Status      *string `gorm:"not null;check:status IN ('published', 'unpublished')"`
		Featured    bool    `gorm:"default:false"`
		Stock       uint    `gorm:"-"`
		IsChild     bool    `gorm:"default:false"`
		ParentID    *uint
		Color       string
		Size        string
		BrandID     *uint
		Attributes  []Attributes
		Images      []models.ProductImage `gorm:"foreignKey:ProductID"`
	}

	if err := c.BindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := config.DB.Begin()
	parent := models.Product{
		Name:        payload.Name,
		Description: payload.Description,
		SKU:         payload.SKU,
		Barcode:     payload.Barcode,
		Price:       payload.Price,
		Currency:    payload.Currency,
		CategoryID:  payload.CategoryID,
		BrandID:     payload.BrandID,
		Status:      payload.Status,
		Featured:    payload.Featured,
		Color:       payload.Color,
		Size:        payload.Size,
		Images:      payload.Images,
		Inventory: &models.Inventory{
			StockLevel: int(payload.Stock),
			InOpen:     0,
			ChangeType: "restock",
			ChangeDate: time.Now(),
		},
	}

	if err := tx.Create(&parent).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	if len(payload.Attributes) != 0 {
		var variations []models.Product
		var images []models.ProductImage

		for _, attribute := range payload.Attributes {

			for _, variation := range attribute.Variation {

				variations = append(variations, models.Product{
					Name:        parent.Name,
					Description: parent.Description,
					SKU:         parent.SKU + "-" + variation.Size + "-" + attribute.Color,
					Barcode:     parent.Barcode,
					Price:       parent.Price,
					Currency:    parent.Currency,
					CategoryID:  parent.CategoryID,
					Status:      parent.Status,
					IsChild:     true,
					ParentID:    &parent.ID,
					Color:       attribute.Color,
					Size:        variation.Size,
					BrandID:     parent.BrandID,
					Inventory: &models.Inventory{
						StockLevel: variation.Quantity,
						InOpen:     0,
						ChangeType: "restock",
						ChangeDate: time.Now(),
					},
				})

			}
			images = append(images, models.ProductImage{
				ProductID: parent.ID,
				Image:     attribute.Image,
				Color:     &attribute.Color,
			})

		}

		if err := tx.Create(&variations).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create variations", "error": err.Error()})
			return
		}
		if err := tx.Create(&images).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create Variation images", "error": err.Error()})
			return
		}

	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{"message": "Product added successfully"})
}

// GetProducts retrieves all products with their category and reviews
func SearchProducts(c *gin.Context) {
	var params utils.Parameters
	if c.Bind(&params) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind identifier parameters."})
		return
	}

	searchQuery := ""

	if params.Key != "" {
		searchQuery = "products.name ILIKE '%" + params.Key + "%' OR sku ILIKE '%" + params.Key + "%' OR categories.name ILIKE '%" + params.Key + "%' OR brands.name ILIKE '%" + params.Key + "%' OR size ILIKE '%" + params.Key + "%'"
	}

	type Product struct {
		ID   uint   `gorm:"primarykey"`
		Name string `gorm:"size:150;not null"`
		// Description  string          `gorm:"type:text"`
	}

	var products []*Product

	config.DB.Model(&products).
		Select(`products.name, products.id`).
		Joins("LEFT JOIN categories ON products.category_id = categories.id").
		Joins("LEFT JOIN brands ON products.brand_id = brands.id").
		Where(searchQuery).
		Where("products.status = ? AND products.parent_id IS NULL", "published").
		Group("products.id").Find(&products)

	c.JSON(http.StatusOK, &products)
}
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
	type Brand struct {
		ID   uint        `gorm:"primarykey"`
		Name null.String `gorm:"size:100;not null"`
	}
	type Product struct {
		gorm.Model
		Name            string  `gorm:"size:150;not null"`
		Description     string  `gorm:"type:text"`
		SKU             string  `gorm:"size:150;not null;unique;index"`
		Barcode         *string `gorm:"size:150"`
		Price           float64 `gorm:"type:decimal(10,2);not null"`
		Currency        string  `gorm:"size:3; not null"`
		BrandID         *uint
		Brand           Brand           `gorm:"foreignKey:BrandID"`
		CategoryID      uint            `gorm:"not null"`
		Category        models.Category `gorm:"foreignKey:CategoryID"`
		Featured        bool            `gorm:"default:false"`
		Status          *string         `gorm:"not null;check:status IN ('published', 'unpublished')"`
		Inventory       *Inventory      `gorm:"foreignKey:ProductID"`
		InventoryStatus bool            `gorm:"column:inventory_status"`
		TotalReviews    int
		Rating          int
		Images          []models.ProductImage `gorm:"foreignKey:ProductID"`
	}

	var products []*Product
	var model *gorm.DB

	model = config.DB.Debug().Model(&products).Preload("Category").Preload("Inventory").Preload("Brand").Preload("Images").
		Select(`products.*, 
				count(reviews.id) as total_reviews,
				AVG(reviews.rating)::int as rating,
				CASE 
					WHEN EXISTS (
						SELECT 1 
						FROM products AS child
						JOIN inventories ON inventories.product_id = child.id
						WHERE child.parent_id = products.id 
						AND inventories.stock_level > 0
					) THEN true
					WHEN EXISTS (
						SELECT 1 
						FROM inventories 
						WHERE inventories.product_id = products.id 
						AND inventories.stock_level > 0
					) THEN true
					ELSE false
				END AS inventory_status

			`).
		Joins("LEFT JOIN reviews ON products.id = reviews.product_id").
		Where(querystring).
		Where("is_child = ?", false).
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
	type Brand struct {
		ID   uint        `gorm:"primarykey"`
		Name null.String `gorm:"size:100;not null"`
	}
	type Product struct {
		gorm.Model
		Name            string                `gorm:"size:150;not null"`
		Description     string                `gorm:"type:text"`
		SKU             string                `gorm:"size:150;not null;unique;index"`
		Barcode         *string               `gorm:"size:150"`
		Price           float64               `gorm:"type:decimal(10,2);not null"`
		Currency        string                `gorm:"size:3; not null"`
		Images          []models.ProductImage `gorm:"foreignKey:ProductID"`
		BrandID         *uint
		Brand           Brand           `gorm:"foreignKey:BrandID"`
		CategoryID      uint            `gorm:"not null"`
		Category        models.Category `gorm:"foreignKey:CategoryID"`
		Status          *string         `gorm:"not null;check:status IN ('published', 'unpublished')"`
		Inventory       *Inventory      `gorm:"foreignKey:ProductID"`
		InventoryStatus bool            `gorm:"column:inventory_status"`
		TotalReviews    int
		Rating          int
	}

	var products []*Product
	var model *gorm.DB

	model = config.DB.Model(&products).Preload("Category").Preload("Inventory").Preload("Brand").Preload("Images").
		Select(`products.*, 
				count(reviews.id) as total_reviews,
				AVG(reviews.rating)::int as rating,
				CASE 
					WHEN EXISTS (
						SELECT 1 
						FROM products AS child
						JOIN inventories ON inventories.product_id = child.id
						WHERE child.parent_id = products.id 
						AND inventories.stock_level > 0
					) THEN true
					WHEN EXISTS (
						SELECT 1 
						FROM inventories 
						WHERE inventories.product_id = products.id 
						AND inventories.stock_level > 0
					) THEN true
					ELSE false
				END AS inventory_status
			`).
		Joins("LEFT JOIN reviews ON products.id = reviews.product_id").
		Where(querstring).
		Where("is_child = false").
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
	type Brand struct {
		ID   uint        `gorm:"primarykey"`
		Name null.String `gorm:"size:100;not null"`
	}
	type Product struct {
		gorm.Model
		Name            string                `gorm:"size:150;not null"`
		Description     string                `gorm:"type:text"`
		SKU             string                `gorm:"size:150;not null;unique;index"`
		Barcode         *string               `gorm:"size:150"`
		Price           float64               `gorm:"type:decimal(10,2);not null"`
		Currency        string                `gorm:"size:3; not null"`
		Images          []models.ProductImage `gorm:"foreignKey:ProductID"`
		BrandID         *uint
		Brand           Brand           `gorm:"foreignKey:BrandID"`
		CategoryID      uint            `gorm:"not null"`
		Category        models.Category `gorm:"foreignKey:CategoryID"`
		Status          *string         `gorm:"not null;check:status IN ('published', 'unpublished')"`
		Inventory       *Inventory      `gorm:"foreignKey:ProductID"`
		InventoryStatus bool            `gorm:"column:inventory_status"`
		TotalReviews    int
		Rating          int
	}

	var products []*Product
	var model *gorm.DB

	model = config.DB.Model(&products).Preload("Category").Preload("Inventory").Preload("Brand").Preload("Images").
		Select(`products.*, 
				count(reviews.id) as total_reviews,
				AVG(reviews.rating)::int as rating,
				CASE 
					WHEN EXISTS (
						SELECT 1 
						FROM products AS child
						JOIN inventories ON inventories.product_id = child.id
						WHERE child.parent_id = products.id 
						AND inventories.stock_level > 0
					) THEN true
					WHEN EXISTS (
						SELECT 1 
						FROM inventories 
						WHERE inventories.product_id = products.id 
						AND invenstock_level > 0
					) THEN true
					ELSE false
				END AS inventory_status
			`).
		Joins("LEFT JOIN reviews ON products.id = reviews.product_id").
		Joins("LEFT JOIN order_items on products.id = order_items.product_id").
		Where(querstring).
		Where("is_child = ?", false).
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
func GetSingleProduct(c *gin.Context) {
	productID := c.Param("id")

	type Category struct {
		CategoryID       *uint `gorm:"-"`
		SubCategoryID    *uint `gorm:"-"`
		SubSubCategoryID *uint `gorm:"-"`
	}

	type Brand struct {
		ID   uint        `gorm:"primarykey"`
		Name null.String `gorm:"size:100;not null"`
	}
	type Product struct {
		gorm.Model
		Name         string   `gorm:"size:150;not null"`
		Description  string   `gorm:"type:text"`
		SKU          string   `gorm:"size:150;not null;unique;index"`
		Barcode      *string  `gorm:"size:150"`
		Price        float64  `gorm:"type:decimal(10,2);not null"`
		Currency     string   `gorm:"size:3; not null"`
		CategoryID   uint     `gorm:"not null"`
		Category     Category `gorm:"foreignKey:CategoryID"`
		BrandID      *uint
		Brand        Brand   `gorm:"foreignKey:BrandID"`
		Status       *string `gorm:"not null;check:status IN ('published', 'unpublished')"`
		Inventory    int     `gorm:"column:stock_level"`
		TotalReviews int
		Rating       int
		Variation    json.RawMessage
		Path         pq.Int64Array `gorm:"column:path" json:"-"`
		Images       []models.ProductImage
	}

	var product *Product

	model := config.DB.Model(&models.Product{}).Preload("Category").Preload("Brand").
		Select(`products.*, 
				count(reviews.id) as total_reviews,
				cat_path.path,
				AVG(reviews.rating)::int as rating,
				COALESCE(
					json_agg(
						json_build_object(
						'id', variations.id,
						'size', variations.size,
						'color', variations.color

						)
					)FILTER (WHERE variations.id IS NOT NULL),
            		'[]'
				) AS variation,
				inventories.stock_level as stock_level
			`).
		Joins("LEFT JOIN reviews ON products.id = reviews.product_id").
		Joins("LEFT JOIN inventories ON products.id = inventories.product_id").
		Joins("LEFT JOIN products AS variations ON variations.parent_id = products.id AND variations.is_child = true").
		Joins(`
			LEFT JOIN (WITH RECURSIVE category_hierarchy AS (
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
			ORDER BY path) as cat_path on cat_path.id = products.category_id
		`).
		Where("products.id = ?", productID).
		Group("products.id, cat_path.path, inventories.stock_level").
		First(&product)

	if model.Error != nil {
		if errors.Is(model.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "No product found"})
			return

		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": model.Error.Error()})
		return
	}

	config.DB.Model([]models.ProductImage{}).Where("product_id = ?", productID).Find(&product.Images)
	for i, v := range product.Path {
		if i == 0 {
			categoryID := uint(v)
			product.Category.CategoryID = &categoryID
		}
		if i == 1 {
			categoryID := uint(v)
			product.Category.SubCategoryID = &categoryID
		}
		if i == 2 {
			categoryID := uint(v)
			product.Category.SubSubCategoryID = &categoryID
		}
	}

	c.JSON(http.StatusOK, &product)
}

// GetProduct retrieves a single product by its ID
func GetSingleProductV2(c *gin.Context) {
	productID := c.Param("id")

	type AttributeVariation struct {
		ID       uint   `gorm:"primarykey"`
		Size     string `json:"Size"`
		Quantity int    `json:"Quantity"`
	}
	type ColorMap struct {
		Image     string               `json:"Image"`
		Variation []AttributeVariation `json:"variation"`
	}

	type ColorGroup struct {
		Color     string               `json:"Color"`
		Image     string               `json:"Image"`
		Variation []AttributeVariation `json:"variation"`
	}

	type Variation struct {
		ID        uint              `gorm:"primarykey"`
		SKU       string            `gorm:"size:150;not null;unique;index"`
		Inventory *models.Inventory `gorm:"foreignKey:ProductID;references:ID"`
		Color     string
		Size      string
		Images    []models.ProductImage `gorm:"foreignKey:ProductID;references:ID"`
	}

	type Product struct {
		gorm.Model
		Name        string          `gorm:"size:150;not null"`
		Description string          `gorm:"type:text"`
		SKU         string          `gorm:"size:150;not null;unique;index"`
		Barcode     *string         `gorm:"size:150"`
		Price       float64         `gorm:"type:decimal(10,2);not null"`
		Currency    string          `gorm:"size:3; not null"`
		CategoryID  uint            `gorm:"not null"`
		Category    models.Category `gorm:"foreignKey:CategoryID"`
		Status      *string         `gorm:"not null;check:status IN ('published', 'unpublished')"`
		Featured    bool            `gorm:"default:false"`
		Stock       uint            `gorm:"-"`
		IsChild     bool            `gorm:"default:false"`
		ParentID    *uint
		Color       string
		Size        string
		BrandID     *uint
		Brand       models.Brand          `gorm:"foreignKey:BrandID;refrences:BrandID"`
		Images      []models.ProductImage `gorm:"foreignKey:ProductID"`
		Variations  []Variation           `gorm:"-" json:"-"`
		Attributes  []ColorGroup          `gorm:"-"`
	}

	var product *Product
	// var variations []Variation

	model := config.DB.Model(&models.Product{}).Preload("Category").Preload("Brand").Preload("Images").First(&product, productID)

	if model.Error != nil {
		if errors.Is(model.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "No product found"})
			return

		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": model.Error.Error()})
		return
	}
	config.DB.Model([]models.Product{}).Preload("Images").Preload("Inventory").Where("parent_id = ? AND is_child = true", productID).Find(&product.Variations)

	config.DB.Model([]models.ProductImage{}).Where("product_id = ?", productID).Find(&product.Images)

	colorMap := make(map[string]ColorMap)

	for _, product := range product.Variations {
		variation := AttributeVariation{
			ID:       product.ID,
			Size:     product.Size,
			Quantity: product.Inventory.StockLevel,
		}
		// Get the current value (or zero value if not present)
		colorGroup := colorMap[product.Color]

		// Append the variation
		colorGroup.Variation = append(colorGroup.Variation, variation)

		// Put it back into the map
		colorMap[product.Color] = colorGroup
	}

	for _, img := range product.Images {
		// Get the current value (or zero value if not present)
		colorGroup := colorMap[*img.Color]

		// Append the variation
		colorGroup.Image = img.Image

		// Put it back into the map
		colorMap[*img.Color] = colorGroup
	}

	var output []ColorGroup

	for color, obj := range colorMap {
		output = append(output, ColorGroup{
			Color:     color,
			Image:     obj.Image,
			Variation: obj.Variation,
		})
	}
	product.Attributes = output

	c.JSON(http.StatusOK, &product)
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
