package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(router *gin.Engine) {
	orders := router.Group("/orders")
	{
		orders.POST("/", controllers.CreateOrder)                             // Create an order
		orders.GET("/:id", controllers.GetOrder)                              // Get an order by ID
		orders.POST("/shippping-address/", controllers.CreateShippingAddress) // Get an order by ID
	}
}
