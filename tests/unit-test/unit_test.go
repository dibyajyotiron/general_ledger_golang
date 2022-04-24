package unit_test

import (
	"strings"
	"testing"

	asrt "github.com/stretchr/testify/assert"

	"general_ledger_golang/models"
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
