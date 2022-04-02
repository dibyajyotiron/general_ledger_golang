package v1

import (
	"encoding/json"
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
	reqBody := util.ParseReqBodyToMap(c.Request.Body)

	opType := reqBody["type"]
	memo := reqBody["memo"]
	entries := reqBody["entries"]
	metadata := reqBody["metadata"]

	x, _ := json.Marshal(reqBody)
	fmt.Printf("Req Body: %+v", string(x))
	fmt.Printf("Req Body Metadata Type: %+v", reflect.TypeOf(entries))

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
		appGin.Response(http.StatusInternalServerError, e.ERROR, map[string]interface{}{"message": "Creating operation resulted in error!"})
		return
	}

	// return the operation
	appGin.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{"operation": foundOp})
	return
}
