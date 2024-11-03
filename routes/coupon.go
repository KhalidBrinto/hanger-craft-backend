package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func CuponRoutes(router *gin.Engine) {
	coupon := router.Group("/api/coupons")
	{
		coupon.POST("/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.CreateCoupon)
		coupon.GET("", controllers.GetCoupons)
		coupon.PUT("/:id/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.UpdateCoupon)
		coupon.DELETE("/:id/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.DeleteCoupon)
	}
}
