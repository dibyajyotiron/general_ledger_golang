package routers

import (
	"general_ledger_golang/middleware"
	v1 "general_ledger_golang/routers/api/v1"

	"github.com/gin-gonic/gin"

	_ "general_ledger_golang/docs"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()

	// default middlewares
	r.Use(gin.Logger())

	// Custom middlewares
	r.Use(middleware.CORS())
	r.Use(gin.CustomRecovery(middleware.ErrorHandler))

	// apiV1 groups
	apiV1 := r.Group("/api/v1")

	// Jwt unprotected routes
	apiV1.GET("/test", v1.TestAppStatus)

	// Books route
	apiV1Groups := apiV1.Group("/books")
	apiV1Groups.POST("/", middleware.UseRequestBody(), v1.CreateOrUpdateBook)
	apiV1Groups.GET("/:bookId", v1.GetBook)
	// Jwt protected routes

	apiV1.GET("/secured/test", middleware.JWT(), v1.TestAppStatus)

	// apiV2 Jwt protected routes
	apiV2 := r.Group("/api/v2", middleware.JWT())
	{
		apiV2.GET("/test", v1.TestAppStatus)
	}

	return r
}
