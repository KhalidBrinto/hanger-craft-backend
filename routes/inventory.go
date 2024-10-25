package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func InventoryRoutes(router *gin.Engine) {
	inventory := router.Group("/api/inventory")
	{
		inventory.POST("/restock/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.RestockProduct) // Add stock (restock)
		inventory.GET("", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.GetInventory)             // Add stock (restock)
	}
}
