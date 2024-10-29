package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(router *gin.Engine) {
	categories := router.Group("/api/categories")
	{
		categories.POST("/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.CreateCategory)
		categories.GET("", controllers.GetCategories)
		categories.GET("/sub-category/:parent_id", controllers.GetSubCategories)
		categories.GET("/:id", controllers.GetCategory)
		categories.PUT("/:id/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.UpdateCategory)
		categories.DELETE("/:id/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.DeleteCategory)
	}
}
