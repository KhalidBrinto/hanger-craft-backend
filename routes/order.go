package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(router *gin.Engine) {
	orders := router.Group("/api/orders")
	{
		orders.POST("/", middlewares.AuthMiddleware(), controllers.CreateOrder)
		orders.GET("/:id", middlewares.AuthMiddleware(), controllers.GetOrderByID)
		orders.GET("", middlewares.AuthMiddleware(), controllers.GetOrders)
		orders.PUT("/dispatch/:id/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.DispatchOrder)
		orders.PUT("/cancel/:id/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.CancelOrder)
		orders.PUT("/:id/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.UpdateOrderStatus)
	}
}
