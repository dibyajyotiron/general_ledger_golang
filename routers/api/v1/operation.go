package v1

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"general_ledger_golang/pkg/app"
	"general_ledger_golang/pkg/e"
	"general_ledger_golang/pkg/util"
	"general_ledger_golang/service/operation_service"
)

func PostOperation(c *gin.Context) {
	appGin := app.Gin{C: c}
	reqBody := util.GetReqBodyFromCtx(c)

	opType := reqBody["type"]
	memo := reqBody["memo"]
	entries := reqBody["entries"]
	metadata := reqBody["metadata"]

	fmt.Printf("\nReq Body: %+v", reqBody)
	fmt.Printf("\nReq Body Metadata Type: %+v", reflect.TypeOf(entries))

	opMap := map[string]interface{}{
		"type":     opType,
		"memo":     memo,
		"entries":  entries,
		"metadata": metadata,
	}

	opService := &operation_service.OperationService{}
	foundOp, err := opService.PostOperation(opMap)

	if err != nil || foundOp == nil {
		fmt.Printf("Creating Operation Failed, error: %+v", err)
		appGin.Response(http.StatusInternalServerError, e.ERROR, map[string]interface{}{
			"message": "Creating operation resulted in error!",
			"error":   err.Error(),
		})
		return
	}

	// return the operation
	appGin.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{"operation": foundOp})
	return
}
