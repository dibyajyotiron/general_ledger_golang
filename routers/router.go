package routers

import (
	v1 "general_ledger_golang/routers/api/v1"

	"github.com/gin-gonic/gin"

	_ "general_ledger_golang/docs"

	middleware "general_ledger_golang/middleware/cors"
	error_middleware "general_ledger_golang/middleware/errors"
	"general_ledger_golang/middleware/jwt"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()

	// default middlewares
	r.Use(gin.Logger())

	// Custom middlewares
	r.Use(middleware.CORS())
	r.Use(gin.CustomRecovery(error_middleware.Error_Handler))

	// apiV1 groups
	apiV1 := r.Group("/api/v1")

	// Jwt unprotected routes
	apiV1.GET("/test", v1.TestAppStatus)

	// Jwt protected routes
	apiV1.Use(jwt.JWT())
	{
		apiV1.GET("/secured/test", v1.TestAppStatus)
	}

	// apiV2 Jwt protected routes
	apiV2 := r.Group("/api/v2", jwt.JWT())
	{
		apiV2.GET("/test", v1.TestAppStatus)
	}

	return r
}
