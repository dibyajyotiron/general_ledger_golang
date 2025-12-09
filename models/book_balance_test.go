package models

import (
	"strings"
	"testing"

	"general_ledger_golang/domain"
)

func TestMain(t *testing.T) {
	// GenerateUpsertCteQuery will generate multiple upsert queries

	queries, params, err := GenerateUpsertCteQuery([]domain.OperationEntry{
		{AssetId: "btc", BookId: "3", Value: "-1"},
		{AssetId: "btc", BookId: "4", Value: "1"},
	}, map[string]interface{}{
		"operation": "BLOCK",
	})

	if len(queries) != 2 {
		t.Fatalf("Queries should have 2 elements")
	}

	if len(params) != 2 {
		t.Fatalf("Params should have 2 elements")
	}

	if err != nil {
		t.Fatalf("Err should be nil")
	}

	t.Log(strings.Join(queries, "\n"), params, err)
}
