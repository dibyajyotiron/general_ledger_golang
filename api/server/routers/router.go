package routers

import (
	v12 "general_ledger_golang/api/server/routers/api/v1"
	"general_ledger_golang/middleware"
	"general_ledger_golang/models"

	"github.com/gin-gonic/gin"
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
	apiV1.GET("/test", v12.TestAppStatus)

	// Books route
	apiV1BooksGroup := apiV1.Group("/books")
	apiV1BooksGroup.POST("/", middleware.UseRequestBody(), v12.CreateOrUpdateBook)
	apiV1BooksGroup.GET("/:bookId", v12.GetBook)
	apiV1BooksGroup.GET("/:bookId/balance", v12.GetBookBalance)

	// Operations route
	apiV1OperationsGroup := apiV1.Group("/operations")
	apiV1OperationsGroup.POST("/", middleware.UseRequestBody(), middleware.ReqBodySanitizer(models.ValidatePostOperation), v12.PostOperation)
	apiV1OperationsGroup.GET("/", v12.GetOperationByMemo)
	// Jwt protected routes

	apiV1.GET("/secured/test", middleware.JWT(), v12.TestAppStatus)

	// apiV2 Jwt protected routes
	apiV2 := r.Group("/api/v2", middleware.JWT())
	{
		apiV2.GET("/test", v12.TestAppStatus)
	}

	return r
}
