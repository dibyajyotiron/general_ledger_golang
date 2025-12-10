package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"general_ledger_golang/internal/app"
	"general_ledger_golang/internal/e"
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
		return
	}
	appG.Response(
		http.StatusInternalServerError,
		e.ERROR,
		"Something Went Wrong...",
	)
	c.Abort()
	return
}
