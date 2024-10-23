package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(router *gin.Engine) {
	products := router.Group("/api/products")
	{
		products.POST("/", controllers.CreateProduct)
		products.GET("", controllers.GetProducts)
		products.GET("/:id", controllers.GetProduct)
		products.PUT("/:id/", controllers.UpdateProduct)
		products.DELETE("/:id/", controllers.DeleteProduct)
	}

	productAttributes := router.Group("/api/product-attributes")
	{
		productAttributes.POST("/", controllers.CreateProductAttribute)
		productAttributes.GET("", controllers.GetProductAttributes)
		productAttributes.PUT("/:id/", controllers.UpdateProductAttribute)
		productAttributes.DELETE("/:id/", controllers.DeleteProductAttribute)
	}
}
