package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func ContentRoutes(router *gin.Engine) {
	content := router.Group("/api/content")
	{
		content.POST("/upload-banner-image/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.AddBannerImages)
		content.GET("/banner-image", controllers.GetBannerImages)
	}
}
