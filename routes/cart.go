package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func CartRoutes(router *gin.Engine) {
	cartRoutes := router.Group("/api/cart")
	{
		cartRoutes.POST("/", controllers.CreateShoppingCart)
		cartRoutes.GET("/user/:user_id", controllers.GetShoppingCartByUserID)
		cartRoutes.POST("/item/", controllers.AddCartItem)
		cartRoutes.PUT("/item/:id", controllers.UpdateCartItem)
		cartRoutes.DELETE("/item/:id/", controllers.RemoveCartItem)
		cartRoutes.DELETE("/:uuid/", controllers.DeleteShoppingCart)
	}
}
