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

	router.GET("/", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, "Hanger Craft API Service health is OK") })
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

	config.ConnectDatabase()
	router.Run(":3000")
}
