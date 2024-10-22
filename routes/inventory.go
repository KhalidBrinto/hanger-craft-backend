package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func InventoryRoutes(router *gin.Engine) {
	inventory := router.Group("/inventory")
	{
		inventory.POST("/restock", controllers.RestockProduct) // Add stock (restock)
	}
}
