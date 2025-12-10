package unit_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	asrt "github.com/stretchr/testify/assert"

	"general_ledger_golang/domain"
	"general_ledger_golang/internal/config"
	"general_ledger_golang/internal/logger"
	"general_ledger_golang/models"
)

func TestQueryGeneration(t *testing.T) {
	assert := asrt.New(t)
	t.Run("Test_Bulk_Query_Generation", func(t *testing.T) {
		entries := []domain.OperationEntry{
			{BookId: "4112314", AssetId: "inr", Value: "20000"},
			{BookId: "4112313", AssetId: "inr", Value: "-20000"},
		}
		logger.Logger.Infof("Entries: %+v", entries)
		metadata := map[string]interface{}{"operation": "DEPOSIT"}
		q, p, e := models.GenerateBulkUpsertQuery(entries, metadata)
		assert.Nil(e)
		assert.NotEqual("", p, "Params slice should not be empty")
		assetParam := p[0].(string)
		bookParam := p[1].(string)
		opParam := p[2].(string)
		valueParam := p[3].(string)
		assert.Equal(len(entries), len(strings.Split(assetParam, ",")))
		assert.Equal(len(entries), len(strings.Split(bookParam, ",")))
		assert.Equal(len(entries), len(strings.Split(opParam, ",")))
		assert.Equal(len(entries), len(strings.Split(valueParam, ",")))
		logger.Logger.Printf("Query: %+v", strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", ""), "\\", ""))
		//logger.Logger.Printf("Params: %+v", p)
	})
	t.Run("Test_Upsert_Query_Generation", func(t *testing.T) {
		entries := []domain.OperationEntry{
			{BookId: "4112314", AssetId: "inr", Value: "20000"},
			{BookId: "4112313", AssetId: "inr", Value: "-20000"},
		}
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

	config.Setup("../../internal/config/")
	c := *config.GetConfig()

	fmt.Printf("Config App: %+v\n", *c.AppSetting)
	fmt.Printf("Config Server: %+v\n", *c.ServerSetting)
	fmt.Printf("Config Redis: %+v\n", *c.RedisSetting)
	fmt.Printf("Config Database: %+v\n", *c.DatabaseSetting)

	assert.Equal(c.AppSetting.JwtSecret, "usx1957-213123123123-12312sa7687-23424")
	assert.Equal(c.RedisSetting.Host, "127.0.0.1:6379")
}
