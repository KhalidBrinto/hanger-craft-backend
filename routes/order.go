package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(router *gin.Engine) {
	orders := router.Group("/api/orders")
	{
		orders.POST("/", middlewares.AuthMiddleware(), controllers.CreateOrder)    // Create an order
		orders.GET("/:id", middlewares.AuthMiddleware(), controllers.GetOrderByID) // Get an order by ID
		orders.GET("", middlewares.AuthMiddleware(), controllers.GetOrders)        // Get an order by ID
	}
}
