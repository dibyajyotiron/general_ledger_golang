package models

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"general_ledger_golang/domain"
	"general_ledger_golang/pkg/logger"
)

type BookBalance struct {
	Model
	BookId        string  `gorm:"primaryKey;index;column:bookId" json:"bookId"`
	AssetId       string  `gorm:"primaryKey;index;column:assetId" json:"assetId"`
	OperationType string  `gorm:"primaryKey;index;column:operationType" json:"operationType"`
	Balance       float64 `gorm:"type:numeric(32,8);check:non_negative_balance,balance >= 0 OR \"operationType\" != 'OVERALL' OR \"bookId\" = '1'" json:"balance"`
}

const (
	OverallOperation string = "OVERALL"
)

func (bB *BookBalance) ModifyBalance(ctx context.Context, entries []domain.OperationEntry, metadata map[string]interface{}, gormDB *gorm.DB) error {
	log := logger.Logger.WithFields(logrus.Fields{
		"memo": fmt.Sprint(metadata["memo"]),
		"op": map[string]interface{}{
			"entries":  entries,
			"metadata": metadata,
		},
	})

	if metadata == nil {
		return errors.New("metadata is not present, creation of book balance depends on metadata[\"operation\"], please send metadata with operation")
	}

	if gormDB == nil {
		gormDB = db
	}

	metaCopy := map[string]interface{}{}
	for k, v := range metadata {
		metaCopy[k] = v
	}
	metaCopy["operation"] = OverallOperation

	opEntries := make([]domain.OperationEntry, len(entries))
	copy(opEntries, entries)
	bB.sortEntries(opEntries)

	queryList, params, err := GenerateUpsertCteQuery(opEntries, metaCopy)
	if err != nil {
		return err
	}

	log.Infof("Executing -> quries: %+v, params: %+v", queryList, params)

	for i, query := range queryList {
		t := gormDB.WithContext(ctx).Debug().Exec(query, params[i]...)
		if t.Error != nil {
			log.WithFields(map[string]interface{}{
				"q": map[string]interface{}{
					"query": strings.ReplaceAll(strings.ReplaceAll(query, "\t", " "), "\n", " "),
					"vars":  params[i],
				},
			}).Errorf("DB error, %+v", t.Error)

			return errors.New(t.Error.Error())
		}
	}

	return nil
}

// GenerateBulkUpsertQuery will generate a single bulkUpsert query
func GenerateBulkUpsertQuery(entries []domain.OperationEntry, metadata map[string]interface{}) (query string, params []interface{}, errs error) {
	var bookIds []string
	var assetIds []string
	var operationTypes []string
	var values []string

	bulkInsertQ := `INSERT INTO book_balances 
				(
					"assetId", 
					"bookId", 
					"operationType", 
					"balance", 
					"createdAt", 
					"updatedAt" 
				) 
			SELECT 
				unnest(string_to_array(?, ',')),
				unnest(string_to_array(?, ',')),
				unnest(string_to_array(?, ',')),
				unnest(string_to_array(?, ',')::numeric[]),
				NOW()::timestamp, 
				NOW()::timestamp `
	updateQ := `UPDATE book_balances 
					SET balance = book_balances.balance + data_table.value::numeric, 
					"updatedAt" = NOW()::timestamp 
				FROM 
					( 
						SELECT unnest(string_to_array(?, ',')) as "assetId",
								unnest(string_to_array(?, ',')) as "bookId", 
								unnest(string_to_array(?, ',')) as "operationType", 
								unnest(string_to_array(?, ',')::numeric[]) as value 
					) as data_table 
				WHERE
					book_balances."assetId" = data_table."assetId" 
					AND book_balances."bookId" = data_table."bookId" 
					AND book_balances."operationType" = data_table."operationType" 
				RETURNING *`

	cteQ := fmt.Sprintf(`
			WITH upsert AS (
                        %s
			)
			%s
			WHERE NOT EXISTS(
				SELECT * FROM upsert
			);
			`, updateQ, bulkInsertQ)
	var errList []string
	for _, entry := range entries {
		if entry.BookId == "1" || entry.BookId == "-1" {
			continue
		}

		operationType, ok := metadata["operation"]
		if !ok {
			return "", nil, errors.New("operation is not present inside metadata, creation of book balance depends on metadata[\"operation\"], please send metadata with operation")
		}

		bookIds = append(bookIds, entry.BookId)
		assetIds = append(assetIds, entry.AssetId)
		values = append(values, entry.Value)
		operationTypes = append(operationTypes, fmt.Sprint(operationType))
	}
	if len(errList) > 0 {
		return "", nil, errors.New(strings.Join(errList, ""))
	}
	params = append(params,
		strings.Join(assetIds, ","),
		strings.Join(bookIds, ","),
		strings.Join(operationTypes, ","),
		strings.Join(values, ","),
		strings.Join(assetIds, ","),
		strings.Join(bookIds, ","),
		strings.Join(operationTypes, ","),
		strings.Join(values, ","),
	)
	return cteQ, params, nil
}

// GenerateUpsertCteQuery will generate multiple upsert queries
func GenerateUpsertCteQuery(entries []domain.OperationEntry, metadata map[string]interface{}) (queryList []string, params [][]interface{}, err error) {
	for _, entry := range entries {
		var paramsSlice []interface{}
		if strings.Contains(os.Getenv("EXCLUDED_BALANCE_BOOK_IDS"), entry.BookId) {
			continue
		}
		operationType, ok := metadata["operation"]
		if !ok {
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
		paramsSlice = append(paramsSlice, gorm.Expr("book_balances.balance + ?::numeric ", entry.Value), entry.AssetId, entry.BookId, fmt.Sprint(operationType))

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
			entry.BookId,
			entry.AssetId,
			fmt.Sprint(operationType),
			entry.Value)

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

func (bB *BookBalance) sortEntries(entries []domain.OperationEntry) {
	sort.SliceStable(entries, func(i, j int) bool {
		iEntry := entries[i]
		jEntry := entries[j]
		sortByBookId := iEntry.BookId < jEntry.BookId
		if iEntry.BookId == jEntry.BookId {
			sortByAssetId := iEntry.AssetId < jEntry.AssetId
			return sortByAssetId
		}
		return sortByBookId
	})
}

// GetBalance fetches the balance
func (bB *BookBalance) GetBalance(ctx context.Context, bookId, assetId, operationType string, tx *gorm.DB) ([]BookBalance, error) {
	var d *gorm.DB
	if bookId == "" {
		return nil, errors.New("BookId is missing")
	}
	if tx != nil {
		d = tx.WithContext(ctx)
	} else {
		d = db.WithContext(ctx)
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

	return balance, nil
}
