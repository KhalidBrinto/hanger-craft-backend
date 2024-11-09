package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.Engine) {
	products := router.Group("/api/products")
	{
		products.GET("/search", controllers.SearchProducts)
		products.POST("/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.CreateProduct)
		products.GET("", controllers.GetProducts)
		products.GET("/:id", controllers.GetSingleProduct)
		products.GET("/new-arrival", controllers.GetNewArrivalProducts)
		products.GET("/trending", controllers.GetTrendingProducts)
		products.PUT("/:id/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.UpdateProduct)
		products.DELETE("/:id/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.DeleteProduct)
	}

	productAttributes := router.Group("/api/product-attributes")
	{
		productAttributes.POST("/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.CreateProductAttribute)
		productAttributes.GET("", controllers.GetProductAttributes)
		productAttributes.PUT("/:id/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.UpdateProductAttribute)
		productAttributes.DELETE("/:id/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.DeleteProductAttribute)
	}
}
