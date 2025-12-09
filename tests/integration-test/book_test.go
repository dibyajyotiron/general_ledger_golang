package integration_test

import (
	"fmt"
	"strings"
	"testing"

	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	asrt "github.com/stretchr/testify/assert"

	_ "general_ledger_golang/pkg/logger"
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
		assert.JSONEq(`{
	"code": 200,
	"msg": "SUCCESS",
	"data": {
		"book": {
			"createdAt": "2022-03-28T01:46:47.877328+05:30",
			"id": 4,
			"metadata": {
				"phone": "+919674753375",
				"user_id": 1
			},
			"name": "dibyajyoti_main_book",
			"updatedAt": "2022-03-28T08:40:57.814863+05:30"
		}
	}
}`,
			strings.Trim(*res, "\n"),
		)
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
