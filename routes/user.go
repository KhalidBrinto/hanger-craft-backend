package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userRoutes := router.Group("/api/user")
	{
		userRoutes.POST("/", controllers.CreateUser)
		// worklogRoutes.GET("/single/:day_identifier", controller.GetWorklogByDayIdentifier)
		// worklogRoutes.GET("/stat", controller.GetWorklogStat)
		// worklogRoutes.POST("/", controller.CreateWorklog)
		// worklogRoutes.PUT("/:uuid/", controller.UpdateWorklog)
		// worklogRoutes.DELETE("/:uuid/", controller.DeleteWorklog)
	}
}
