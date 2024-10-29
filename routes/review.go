package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func ReviewRoutes(router *gin.Engine) {
	reviews := router.Group("/api/reviews")
	{
		reviews.POST("/", middlewares.AuthMiddleware(), controllers.CreateReview)      // Create a new review
		reviews.GET("/:id", controllers.GetReview)                                     // Get a review by ID
		reviews.GET("", controllers.GetCustomerReview)                                 // Get a review by ID
		reviews.GET("/product/:product_id", controllers.GetReviewsByProduct)           // Get all reviews by product ID
		reviews.PATCH("/:id", middlewares.AuthMiddleware(), controllers.UpdateReview)  // Update a review
		reviews.DELETE("/:id", middlewares.AuthMiddleware(), controllers.DeleteReview) // Delete a review
	}
}
