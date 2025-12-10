package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"general_ledger_golang/internal/e"
)

var (
	validateOnce sync.Once
	validate     *validator.Validate
)

func validatorInstance() *validator.Validate {
	validateOnce.Do(func() {
		validate = validator.New()
	})
	return validate
}

// BindAndValidate binds the JSON body into payloadPtr and validates it using go-playground/validator.
func BindAndValidate[T any](key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var payload T
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": e.INVALID_PARAMS,
				"msg":  "Invalid request body",
				"data": err.Error(),
			})
			return
		}

		if err := validatorInstance().Struct(payload); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": e.INVALID_PARAMS,
				"msg":  "Validation failed",
				"data": err.Error(),
			})
			return
		}

		c.Set(key, payload)
		c.Next()
	}
}
