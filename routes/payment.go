package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func PaymentRoutes(router *gin.Engine) {
	payments := router.Group("/payments")
	{
		payments.POST("/", controllers.CreatePayment)                    // Create a payment
		payments.PATCH("/:id/status", controllers.UpdatePaymentStatus)   // Update payment status
		payments.GET("/order/:order_id", controllers.GetPaymentsByOrder) // Get payments by order ID
	}
}
