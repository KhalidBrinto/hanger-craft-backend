package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func PaymentRoutes(router *gin.Engine) {
	payments := router.Group("/payments")
	{
		payments.POST("/", middlewares.AuthMiddleware(), controllers.CreatePayment)                    // Create a payment
		payments.PATCH("/:id/status", middlewares.AuthMiddleware(), controllers.UpdatePaymentStatus)   // Update payment status
		payments.GET("/order/:order_id", middlewares.AuthMiddleware(), controllers.GetPaymentsByOrder) // Get payments by order ID
	}

	paymentOptions := router.Group("/payment-options")
	{
		paymentOptions.POST("/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.AddPaymentOption) // Create a payment
		paymentOptions.PATCH("/:id/", middlewares.AuthMiddleware(), controllers.UpdatePaymentOption)                     // Update payment status
		paymentOptions.GET("", middlewares.AuthMiddleware(), controllers.GetAvailablePaymentOptions)                     // Get payments by order ID
		paymentOptions.GET("/:id", middlewares.AuthMiddleware(), controllers.GetPaymentOptionByID)                       // Get payments by order ID
	}
}
