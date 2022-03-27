package v1

import (
	"encoding/json"
	"fmt"
	"general_ledger_golang/models"
	"general_ledger_golang/pkg/app"
	"general_ledger_golang/pkg/e"
	"general_ledger_golang/pkg/util"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"net/http"
)

func CreateBook(c *gin.Context) {
	appGin := app.Gin{C: c}
	reqBody := util.ParseReqBodyToMap(c.Request.Body)

	if reqBody == nil {
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, map[string]interface{}{"success": false, "error": "Missing Params"})
		return
	}

	name := reqBody["name"].(string)
	metadataBytes, _ := json.Marshal(reqBody["metadata"])
	restrictionsBytes, _ := json.Marshal(reqBody["restrictions"])
	book := models.Book{Name: name, Metadata: datatypes.JSON(metadataBytes), Restrictions: datatypes.JSON(restrictionsBytes)}

	result := book.CreateBook(&book)
	err := result.Error

	if err != nil {
		fmt.Printf("%+v", err)
		appGin.Response(http.StatusInternalServerError, e.INVALID_PARAMS, map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, "info")
	return
}
