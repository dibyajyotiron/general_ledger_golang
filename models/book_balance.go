package models

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"gorm.io/gorm"

	"general_ledger_golang/pkg/logger"
	"general_ledger_golang/pkg/logging"
	"general_ledger_golang/pkg/util"
)

type BookBalance struct {
	Model
	BookId        string  `gorm:"primaryKey;index;column:bookId" json:"bookId"`
	AssetId       string  `gorm:"primaryKey;index;column:assetId" json:"assetId"`
	OperationType string  `gorm:"primaryKey;index;column:operationType" json:"operationType"`
	Balance       float64 `gorm:"type:numeric(32,8) USING balance::numeric;check:non_negative_balance,balance >= 0 OR \"operationType\" != 'OVERALL' OR \"bookId\" = '1'" json:"balance"`
}

const (
	OverallOperation string = "OVERALL"
)

func (bB *BookBalance) ModifyBalance(operation map[string]interface{}) error {
	genericOp := util.DeepCopyMap(operation)

	overallOp := util.DeepCopyMap(operation)
	if _, ok := overallOp["metadata"]; !ok {
		return errors.New("metadata is not present, creation of book balance depends on metadata[\"operation\"], please send metadata with operation")
	}
	overallOp["metadata"].(map[string]interface{})["operation"] = OverallOperation

	operations := []map[string]interface{}{genericOp, overallOp}

	// sort the operation entries in place, everything is reference type, so sorting in-place works fine
	for _, op := range operations {
		logging.Info("Typeof entries: ", reflect.TypeOf(op["entries"]))
		entries := op["entries"].([]interface{})
		metadata := op["metadata"].(map[string]interface{})
		bB.sortEntries(entries)

		// create the queries by looping over the entries
		queryList, params, err := generateUpsertCteQuery(entries, metadata)
		if err != nil {
			return err
		}

		// execute the queries one by one, if any query errors out, roll back
		for _, query := range queryList {
			for _, param := range params {
				t := db.Exec(query, param...)
				if t.Error != nil {
					log := logger.Logger
					log.WithFields(map[string]interface{}{
						"q": map[string]interface{}{
							"query": strings.ReplaceAll(strings.ReplaceAll(query, "\t", " "), "\n", " "),
							"vars":  param,
						},
					}).Errorf("DB error, %+v", t.Error)

					return errors.New(t.Error.Error())
				}
			}
		}
	}

	return nil
}

func generateUpsertCteQuery(entries []interface{}, metadata map[string]interface{}) ([]string, [][]interface{}, error) {
	var queryList []string
	var params [][]interface{}

	for _, entry2 := range entries {
		entry := entry2.(map[string]interface{})
		var paramsSlice []interface{}
		if util.Includes(entry["bookId"], []interface{}{1, -1, "1", "-1"}) {
			continue
		}
		operationType := metadata["operation"]

		if operationType == nil {
			return nil, nil, errors.New("operation is not present inside metadata, creation of book balance depends on metadata[\"operation\"], please send metadata with operation")
		}

		updateQ := `
				UPDATE book_balances
					SET
					balance = ?,
					"updatedAt" = NOW()::timestamp
				WHERE
					"assetId" = ?
					AND "bookId" = ?
					AND "operationType" = ?
				RETURNING *
			`
		paramsSlice = append(paramsSlice, gorm.Expr("book_balances.balance + ?::numeric ", entry["value"]), entry["assetId"], entry["bookId"], operationType)

		insertQ := `INSERT
			INTO book_balances
			(
				"bookId",
				"assetId",
				"operationType",
				balance,
				"createdAt",
				"updatedAt"
			)
			SELECT 
				?,
				?,
				?,
				?,
				NOW()::timestamp, 
				NOW()::timestamp`

		paramsSlice = append(
			paramsSlice,
			entry["bookId"],
			entry["assetId"],
			operationType,
			entry["value"])

		cteQ := fmt.Sprintf(`
			WITH upsert AS (
                        %s
                    	)
			%s
			WHERE NOT EXISTS(
				SELECT * FROM upsert
			);
			`, updateQ, insertQ)

		queryList = append(queryList, cteQ)
		params = append(params, paramsSlice)
	}
	return queryList, params, nil
}

func (bB *BookBalance) sortEntries(entries []interface{}) {
	sort.SliceStable(entries, func(i, j int) bool {
		iEntry := entries[i].(map[string]interface{})
		jEntry := entries[j].(map[string]interface{})
		// i < j means smallest first, largest last
		sortByBookId := iEntry["bookId"].(string) < jEntry["bookId"].(string)
		// first sort by bookId, if bookIds are matching, then sort by assetId, to guarantee order.
		if iEntry["bookId"] == jEntry["bookId"] {
			sortByAssetId := iEntry["assetId"].(string) < jEntry["assetId"].(string)
			return sortByAssetId
		}
		return sortByBookId
	})
}

// GetBalance fetches the balance
func (bB *BookBalance) GetBalance(bookId, assetId, operationType string, tx *gorm.DB) (*[]BookBalance, error) {
	var d *gorm.DB
	if bookId == "" {
		return nil, errors.New("BookId is missing")
	}
	if tx != nil {
		d = tx
	} else {
		d = db
	}
	var balance []BookBalance
	//err := db.Select("id").Where(Auth{Username: username, Password: password}).First(&auth).Error
	query := d.Where(BookBalance{BookId: bookId})

	if assetId != "" {
		query.Where(BookBalance{AssetId: assetId})
	}

	if operationType != "" {
		query.Where(BookBalance{OperationType: operationType})
	} else {
		query.Where(BookBalance{OperationType: OverallOperation})
	}

	t := query.Select("bookId", "assetId", "balance", "operationType").Find(&balance)

	if t.RowsAffected < 1 {
		return nil, nil
	}

	return &balance, nil
}
