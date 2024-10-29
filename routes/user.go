package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	userRoutes := router.Group("/api/user")
	{
		userRoutes.POST("/", controllers.RegisterCustomer)
		userRoutes.POST("/login/", controllers.LoginUser)
		userRoutes.PUT("/", middlewares.AuthMiddleware(), controllers.UpdateUser)
		userRoutes.GET("/customer", middlewares.AuthMiddleware(), controllers.GetCustomers)
		userRoutes.DELETE("/", middlewares.AuthMiddleware(), controllers.DeleteCustomer)
		userRoutes.DELETE("/:id/", middlewares.AuthMiddleware(), middlewares.CheckIfAdmin(), controllers.DeleteUserByID)
		// worklogRoutes.GET("/single/:day_identifier", controller.GetWorklogByDayIdentifier)
		// worklogRoutes.GET("/stat", controller.GetWorklogStat)
		// worklogRoutes.POST("/", controller.CreateWorklog)
		// worklogRoutes.PUT("/:uuid/", controller.UpdateWorklog)
		// worklogRoutes.DELETE("/:uuid/", controller.DeleteWorklog)
	}
}
