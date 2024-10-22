package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(router *gin.Engine) {
	categories := router.Group("/categories")
	{
		categories.POST("/", controllers.CreateCategory)
		categories.GET("/", controllers.GetCategories)
		categories.GET("/:id", controllers.GetCategory)
		categories.PUT("/:id", controllers.UpdateCategory)
		categories.DELETE("/:id", controllers.DeleteCategory)
	}
}
