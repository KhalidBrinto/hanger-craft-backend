package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func AdminDashboardRoutes(router *gin.Engine) {
	adminDashboardRoutes := router.Group("/api/admin-panel/dashboard")
	adminDashboardRoutes.Use(middlewares.AuthMiddleware())
	adminDashboardRoutes.Use(middlewares.CheckIfAdmin())
	{
		adminDashboardRoutes.GET("/stats", controllers.GetStats)
		adminDashboardRoutes.GET("/monthly-sales", controllers.GetMonthlySales)
		adminDashboardRoutes.GET("/yearly-revenue", controllers.GetYearlyRevenue)
	}
}
