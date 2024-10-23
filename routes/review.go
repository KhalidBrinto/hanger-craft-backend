package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func ReviewRoutes(router *gin.Engine) {
	reviews := router.Group("/reviews")
	{
		reviews.POST("/", controllers.CreateReview)                          // Create a new review
		reviews.GET("/:id", controllers.GetReview)                           // Get a review by ID
		reviews.GET("/product/:product_id", controllers.GetReviewsByProduct) // Get all reviews by product ID
		reviews.PATCH("/:id", controllers.UpdateReview)                      // Update a review
		reviews.DELETE("/:id", controllers.DeleteReview)                     // Delete a review
	}
}
