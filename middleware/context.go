package middleware

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func UseRequestBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		var bodyBytes []byte

		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		// Restore the io.ReadCloser to its original state
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		c.Set("requestBodyBytes", bodyBytes)
		c.Next()
		return
	}
}

// Cont
