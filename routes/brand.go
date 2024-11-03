package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func BrandRoutes(router *gin.Engine) {
	brands := router.Group("/api/brands")
	{
		brands.POST("/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.AddBrand)
		brands.GET("", controllers.GetBrands)
		brands.PUT("/:id/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.UpdateBrand)
		brands.DELETE("/:id/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.DeleteBrand)
	}
}
