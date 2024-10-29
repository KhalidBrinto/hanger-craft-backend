package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func CartRoutes(router *gin.Engine) {
	cartRoutes := router.Group("/api/cart")
	{
		cartRoutes.POST("/", middlewares.AuthMiddleware(), controllers.CreateShoppingCart)
		cartRoutes.GET("", middlewares.AuthMiddleware(), controllers.GetShoppingCartByUserID)
		cartRoutes.POST("/item/", middlewares.AuthMiddleware(), controllers.AddCartItem)
		cartRoutes.PUT("/item/:id", middlewares.AuthMiddleware(), controllers.UpdateCartItem)
		cartRoutes.DELETE("/item/:id/", middlewares.AuthMiddleware(), controllers.RemoveCartItem)
		cartRoutes.DELETE("/:uuid/", middlewares.AuthMiddleware(), controllers.DeleteShoppingCart)
	}

	wishlistRoutes := router.Group("/api/wish-list")
	{
		wishlistRoutes.POST("/", middlewares.AuthMiddleware(), controllers.AddWishlistItem)
		wishlistRoutes.GET("", middlewares.AuthMiddleware(), controllers.GetWishlistByUserID)
		wishlistRoutes.DELETE("/item/:id/", middlewares.AuthMiddleware(), controllers.RemoveWishlistItem)
		wishlistRoutes.DELETE("/", middlewares.AuthMiddleware(), controllers.ClearWishlist)
	}
}
