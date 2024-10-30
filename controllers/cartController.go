package controllers

import (
	"backend/config"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateShoppingCart creates a new shopping cart for a user
func CreateShoppingCart(c *gin.Context) {
	var shoppingCart *models.ShoppingCart

	if err := c.BindJSON(&shoppingCart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure UUID is generated
	shoppingCart.UUID = uuid.New()

	// Save the shopping cart to the database
	if err := config.DB.Create(&shoppingCart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "shopping cart created successfully"})
}

// GetShoppingCartByUserID retrieves the shopping cart by user ID and includes its items
func GetShoppingCartByUserID(c *gin.Context) {
	userID := c.GetUint("user_id")
	var shoppingCart *models.ShoppingCart

	// Use Preload to load associated CartItems
	if err := config.DB.Where("user_id = ?", userID).Preload("CartItems").Preload("CartItems.Product").First(&shoppingCart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shopping cart not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, shoppingCart)
}

func GetWishlistByUserID(c *gin.Context) {
	userID := c.GetUint("user_id")
	var wishList []*models.WishList

	// Use Preload to load associated CartItems
	if err := config.DB.Where("user_id = ?", userID).Preload("Product").Find(&wishList).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No wishlist found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, wishList)
}

// DeleteShoppingCart deletes a shopping cart by UUID
func DeleteShoppingCart(c *gin.Context) {
	cartUUID := c.Param("uuid")
	var shoppingCart *models.ShoppingCart

	if err := config.DB.Where("uuid = ?", cartUUID).First(&shoppingCart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Shopping cart not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := config.DB.Delete(&shoppingCart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shopping cart deleted successfully"})
}

func ClearWishlist(c *gin.Context) {
	userID := c.GetUint("user_id")
	var wishList *models.WishList

	if err := config.DB.Where("user_id = ?", userID).Delete(&wishList).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No wishlist found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "wish-list cleared successfully"})
}

// CreateCartItem adds a new item to the cart
func AddCartItem(c *gin.Context) {
	var cartItem *models.CartItem

	if err := c.BindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if UUID is valid
	if cartItem.CartID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Cart ID"})
		return
	}

	// Save CartItem to the database
	if err := config.DB.Create(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cartItem)
}

func AddWishlistItem(c *gin.Context) {
	var wishlistItem struct {
		ProductID uint `gorm:"not null"`
		UserID    uint `gorm:"not null"`
	}

	if err := c.BindJSON(&wishlistItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	wishlistItem.UserID = c.GetUint("user_id")

	// Save CartItem to the database
	if err := config.DB.Create(&models.WishList{
		ProductID: wishlistItem.ProductID,
		UserID:    wishlistItem.UserID,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item added"})
}

// UpdateCartItem updates the quantity of a cart item
func UpdateCartItem(c *gin.Context) {
	var cartItem *models.CartItem
	cartItemID := c.Param("id")

	if err := config.DB.First(&cartItem, cartItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CartItem not found"})
		return
	}

	if err := c.BindJSON(&cartItem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Save(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cartItem)
}

// DeleteCartItem deletes a cart item by ID
func RemoveCartItem(c *gin.Context) {
	cartItemID := c.Param("id")
	var cartItem *models.CartItem

	if err := config.DB.First(&cartItem, cartItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "CartItem not found"})
		return
	}

	if err := config.DB.Delete(&cartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "CartItem removed successfully"})
}

func RemoveWishlistItem(c *gin.Context) {
	itemID := c.Param("id")
	var wishlistItem *models.WishList

	if err := config.DB.First(&wishlistItem, itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Wish-list item not found"})
		return
	}

	if err := config.DB.Delete(&wishlistItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed successfully"})
}
