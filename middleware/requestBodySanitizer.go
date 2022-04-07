package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"general_ledger_golang/pkg/e"
	"general_ledger_golang/pkg/util"
)

func ReqBodySanitizer(validateFunc func(map[string]interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqBody := util.GetReqBodyFromCtx(c)
		validateReqBody := util.DeepCopyMap(reqBody)

		validateFunc(validateReqBody)

		if validateReqBody["valid"] == false {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": e.INVALID_PARAMS,
				"msg":  e.GetMsg(e.INVALID_PARAMS),
				"data": validateReqBody["errors"],
			})
			c.Abort()
			return
		}
		c.Next()
		return
	}
}
