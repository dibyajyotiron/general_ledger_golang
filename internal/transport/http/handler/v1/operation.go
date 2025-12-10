package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"general_ledger_golang/dto"
	"general_ledger_golang/internal/app"
	"general_ledger_golang/internal/e"
	"general_ledger_golang/internal/logger"
	"general_ledger_golang/internal/service/operation"
)

func PostOperation(c *gin.Context) {
	appGin := app.Gin{C: c}
	payloadVal, _ := c.Get("operationPayload")
	payload := payloadVal.(dto.OperationPayload)

	log := logger.Logger.WithFields(logrus.Fields{
		"memo": payload.Memo,
		"op":   payload,
	})

	log.Infof("Request Received")

	opService := operation.NewOperationService(nil, nil, nil, nil, nil)
	foundOp, err := opService.PostOperation(c.Request.Context(), payload)

	if err != nil || foundOp == nil {
		log.Errorf("Creating Operation Failed, error: %+v", err)
		appGin.Response(http.StatusInternalServerError, e.ERROR, map[string]interface{}{
			"message": "Creating operation resulted in error!",
			"error":   err.Error(),
		})
		return
	}

	// return the operation
	appGin.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{"operation": foundOp})
}

func GetOperationByMemo(c *gin.Context) {
	appGin := app.Gin{C: c}

	memo := c.Query("memo")

	if memo == "" {
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, map[string]interface{}{
			"message": "Memo is not provided!",
		})
		return
	}

	opService := operation.NewOperationService(nil, nil, nil, nil, nil)
	foundOp, err := opService.GetOperation(c.Request.Context(), memo)

	if err != nil {
		logger.Logger.Errorf("Fetching Operation Failed, error: %+v", err)
		appGin.Response(http.StatusInternalServerError, e.ERROR, map[string]interface{}{
			"message": "Fetching operation resulted in error!",
		})
		return
	}

	status := e.SUCCESS
	httpStatus := http.StatusOK

	if foundOp == nil {
		status = e.NOT_EXIST
		httpStatus = http.StatusNotFound
	}
	// return the operation
	appGin.Response(httpStatus, status, map[string]interface{}{"operation": foundOp})
}
