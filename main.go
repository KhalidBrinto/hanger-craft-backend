package main

import (
	"backend/config"
	"backend/middlewares"
	"backend/routes"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.SetTrustedProxies(nil)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	config.ConnectDatabase()

	router.GET("/", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, "Hanger Craft API Service health is OK") })
	// Liveness Probe: Returns 200 if the app is running
	router.GET("/health/liveness", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusOK)
	})

	// Readiness Probe: Checks if the app is ready (e.g., if the DB connection is available)
	router.GET("/health/readiness", func(c *gin.Context) {
		var result int
		err := config.DB.Raw("SELECT 1").Scan(&result).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "error": err.Error()})
			return
		}
		c.AbortWithStatus(http.StatusOK)
	})
	router.Use(middlewares.CORSMiddleware())
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	routes.CartRoutes(router)
	routes.CategoryRoutes(router)
	routes.InventoryRoutes(router)
	routes.OrderRoutes(router)
	routes.PaymentRoutes(router)
	routes.ProductRoutes(router)
	routes.ReviewRoutes(router)
	routes.UserRoutes(router)
	routes.AdminDashboardRoutes(router)

	router.Run(":3000")
}
