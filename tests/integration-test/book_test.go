package integration_test

import (
	"encoding/json"
	"fmt"
	"testing"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	asrt "github.com/stretchr/testify/assert"

	_ "general_ledger_golang/internal/logger"
	"general_ledger_golang/tests"
)

var serverPort = 3000

func TestGetBook(t *testing.T) {
	assert := asrt.New(t)
	t.Run("Test_Book_Without_Balance", func(t *testing.T) {
		rB := tests.TestRequestBuilder{
			URL:    fmt.Sprintf("http://localhost:%d/api/v1/books/%d?balance=false", serverPort, 4),
			BODY:   "",
			METHOD: "GET",
		}
		res, e := rB.MakeApiCall()
		assert.Nil(e)

		var decoded struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
			Data struct {
				Book struct {
					ID       float64                `json:"id"`
					Name     string                 `json:"name"`
					Metadata map[string]interface{} `json:"metadata"`
				} `json:"book"`
			} `json:"data"`
		}

		err := json.Unmarshal([]byte(*res), &decoded)
		assert.Nil(err)
		assert.Equal(200, decoded.Code)
		assert.Equal("SUCCESS", decoded.Msg)
		assert.Equal(float64(4), decoded.Data.Book.ID)
		assert.Equal("dibyajyoti_main_book", decoded.Data.Book.Name)
		assert.Equal("+919674753375", decoded.Data.Book.Metadata["phone"])
		assert.Equal(float64(1), decoded.Data.Book.Metadata["user_id"])
	})
}

//func TestDeposit(t *testing.T) {
//	assert := asrt.New(t)
//
//	opType := "DEPOSIT"
//	assetId := "inr"
//	amount := 10000.7728
//	bookId := 4
//
//	rB := tests.TestRequestBuilder{}
//
//	bodyMap := map[string]interface{}{
//		"type": "TRANSFER",
//		"memo": strconv.Itoa(time.Time{}.Nanosecond()) + "_" + opType,
//		"entries": SM{
//			{
//				"bookId":  "1",
//				"assetId": assetId,
//				"value":   strconv.FormatFloat(-1*amount, 'F', -1, 64),
//			},
//			{
//				"bookId":  bookId,
//				"assetId": assetId,
//				"value":   strconv.FormatFloat(1*amount, 'F', -1, 64),
//			}},
//		"metadata": map[string]interface{}{"operation": "DEPOSIT"},
//	}
//	strBodyMap := rB.MapToJSONString(bodyMap)
//
//	rB.BODY = strBodyMap
//	rB.URL = fmt.Sprintf("http://localhost:%d/operations", serverPort)
//	rB.METHOD = "POST"
//
//	res, _ := rB.MakeApiCall()
//	assert.Equal(200, res.StatusCode)
//	// Now make api call to get balance and check for balance to be same with deposited amount.
//}
