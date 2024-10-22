package main

import (
	"backend/config"
	"backend/routes"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.SetTrustedProxies(nil)
	router.Use(gin.Recovery())

	router.GET("/", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, "Hanger Craft API Service health is OK") })
	// router.Use(middleware.CORSMiddleware())
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	// router.Use(middleware.TokenAuthMiddleware())
	// router.Use(middleware.LoggerMiddleware(logger))

	routes.UserRoutes(router)
	// routes.TemplateRoute(router)
	// routes.DeviceRoute(router)
	// routes.AttendanceRoute(router)
	// routes.SalaryRoute(router)
	// routes.DeleteRoute(router)
	// routes.NotificationRoute(router)
	// routes.DocumentRoute(router)
	// routes.WorklogRoutes(router)
	// routes.BranchPreferenceRoute(router)
	config.ConnectDatabase()
	router.Run(":3000")
}
