package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"general_ledger_golang/dto"
	"general_ledger_golang/internal/app"
	"general_ledger_golang/internal/e"
	"general_ledger_golang/internal/logger"
	"general_ledger_golang/internal/service/book"
)

func GetBook(c *gin.Context) {
	appGin := app.Gin{C: c}
	bookId := c.Param("bookId")
	balanceFetchStr := c.Query("balance")
	balanceFetch, err := strconv.ParseBool(balanceFetchStr)

	if err != nil {
		logger.Logger.Errorf("Parsing of `balance` failed, error: %+v", err)
	}

	bookService := book.NewBookService(nil, nil, nil)
	result, err := bookService.GetBook(c.Request.Context(), bookId, balanceFetch)

	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR, map[string]interface{}{"error": err.Error()})
		return
	}
	if result == nil {
		appGin.Response(http.StatusNotFound, e.NOT_EXIST, map[string]interface{}{"book": result})
		return
	}
	resp := map[string]any{
		"book": result.Book,
	}
	if result.Balances != nil {
		resp["balances"] = result.Balances
	}
	appGin.Response(http.StatusOK, e.SUCCESS, resp)
	return
}

func GetBookBalance(c *gin.Context) {
	appGin := app.Gin{C: c}
	bookId := c.Param("bookId")

	assetId := c.Query("assetId")
	operationType := c.Query("operationType")

	bookService := book.NewBookService(nil, nil, nil)

	result, err := bookService.GetBalance(c.Request.Context(), bookId, assetId, operationType, nil)

	if err != nil {
		appGin.Response(http.StatusInternalServerError, e.ERROR, map[string]interface{}{"error": err.Error()})
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, result)
	return
}

func CreateOrUpdateBook(c *gin.Context) {
	appGin := app.Gin{C: c}
	payloadVal, _ := c.Get("bookPayload")
	payload := payloadVal.(dto.BookPayload)

	bookService := book.NewBookService(nil, nil, nil)
	result, operation, err := bookService.UpsertBook(c.Request.Context(), payload)

	if err != nil {
		fmt.Printf("Book creation failed: %+v", err)
		appGin.Response(http.StatusInternalServerError, e.INVALID_PARAMS, map[string]interface{}{"error": err.Error()})
		return
	}

	appGin.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"book":    result,
		"message": fmt.Sprintf("%v successful", operation),
	})
	return
}
