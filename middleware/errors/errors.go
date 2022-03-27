package error_middleware

import (
	"general_ledger_golang/pkg/app"
	"general_ledger_golang/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler Recovery handler from panics
func ErrorHandler(c *gin.Context, recovered interface{}) {
	appG := app.Gin{C: c}
	if err, ok := recovered.(string); ok {
		appG.Response(
			http.StatusInternalServerError,
			e.ERROR,
			err,
		)
	}
	c.AbortWithStatus(http.StatusInternalServerError)
}
