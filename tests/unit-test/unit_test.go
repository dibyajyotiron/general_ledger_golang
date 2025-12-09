package unit_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	asrt "github.com/stretchr/testify/assert"

	"general_ledger_golang/models"
	"general_ledger_golang/pkg/config"
	"general_ledger_golang/pkg/logger"
)

func TestQueryGeneration(t *testing.T) {
	assert := asrt.New(t)
	t.Run("Test_Bulk_Query_Generation", func(t *testing.T) {
		entries := []interface{}{
			map[string]interface{}{
				"bookId":  "4112314",
				"assetId": "inr",
				"value":   "20000"},
			map[string]interface{}{
				"bookId":  "4112313",
				"assetId": "inr",
				"value":   "-20000",
			}}
		logger.Logger.Infof("Entries: %+v", entries)
		metadata := map[string]interface{}{"operation": "DEPOSIT"}
		q, p, e := models.GenerateBulkUpsertQuery(entries, metadata)
		assert.Nil(e)
		assert.NotEqual("", p, "Params slice should not be empty")
		assert.Len(p[0], len(entries))
		assert.Len(p[1], len(entries))
		assert.Len(p[2], len(entries))
		assert.Len(p[3], len(entries))
		logger.Logger.Printf("Query: %+v", strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", ""), "\\", ""))
		//logger.Logger.Printf("Params: %+v", p)
	})
	t.Run("Test_Upsert_Query_Generation", func(t *testing.T) {
		entries := []interface{}{
			map[string]interface{}{
				"bookId":  "4112314",
				"assetId": "inr",
				"value":   "20000"},
			map[string]interface{}{
				"bookId":  "4112313",
				"assetId": "inr",
				"value":   "-20000",
			}}
		logger.Logger.Infof("Entries: %+v", entries)
		metadata := map[string]interface{}{"operation": "DEPOSIT"}
		q, p, e := models.GenerateUpsertCteQuery(entries, metadata)
		assert.Nil(e)
		assert.NotEqual("", p, "Params slice should not be empty")
		logger.Logger.Printf("Query: %+v", strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(strings.Join(q, "\n"), "\t", ""), "\n", ""), "\\", ""))
		//logger.Logger.Printf("Params: %+v", p)
	})
}

func TestConfig(t *testing.T) {
	assert := asrt.New(t)
	err := godotenv.Load("../../.env")
	assert.Nil(err)

	config.Setup("../../pkg/config/")
	c := *config.GetConfig()

	fmt.Printf("Config App: %+v\n", *c.AppSetting)
	fmt.Printf("Config Server: %+v\n", *c.ServerSetting)
	fmt.Printf("Config Redis: %+v\n", *c.RedisSetting)
	fmt.Printf("Config Database: %+v\n", *c.DatabaseSetting)

	assert.Equal(c.AppSetting.JwtSecret, "usx1957-213123123123-12312sa7687-23424")
	assert.Equal(c.RedisSetting.Host, "127.0.0.1:6379")
}
