package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"

	"general_ledger_golang/models"
	"general_ledger_golang/pkg/app"
	"general_ledger_golang/pkg/e"
	"general_ledger_golang/pkg/util"
)

func GetBook(c *gin.Context) {
	appGin := app.Gin{C: c}
	bookId := c.Param("bookId")
	book := models.Book{}
	result := book.GetBook(bookId)
	if result == nil {
		appGin.Response(http.StatusNotFound, e.NOT_EXIST, map[string]interface{}{"book": result})
		return
	}
	appGin.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{"book": result})
	return
}

func CreateOrUpdateBook(c *gin.Context) {
	appGin := app.Gin{C: c}
	reqBody := util.GetReqBodyFromCtx(c)

	if reqBody == nil {
		appGin.Response(http.StatusBadRequest, e.INVALID_PARAMS, map[string]interface{}{"error": "Missing request body or not a valid json!"})
		return
	}

	name := reqBody["name"].(string)
	metadataBytes, _ := json.Marshal(reqBody["metadata"])
	restrictionsBytes, _ := json.Marshal(reqBody["restrictions"])
	book := models.Book{Name: name, Metadata: datatypes.JSON(metadataBytes), Restrictions: datatypes.JSON(restrictionsBytes)}
	result, operation := book.CreateOrUpdateBook(&book)
	err := result.Error

	if err != nil {
		fmt.Printf("Book creation failed: %+v", err)
		appGin.Response(http.StatusInternalServerError, e.INVALID_PARAMS, map[string]interface{}{"error": err.Error()})
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, fmt.Sprintf("%v successful", operation))
	return
}
