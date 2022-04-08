package book_service

import (
	"errors"
	"strconv"

	"gorm.io/gorm"

	"general_ledger_golang/models"
	"general_ledger_golang/pkg/logger"
	"general_ledger_golang/pkg/util"
)

type BookService struct {
	BookRepository        models.Book
	BookBalanceRepository models.BookBalance
}

func (b *BookService) GetBook(bookId string, withBalance bool) (map[string]interface{}, error) {
	if bookId == "" {
		return nil, errors.New("BookId is empty")
	}
	book, err := b.BookRepository.GetBook(bookId)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, nil
	}

	result := util.StructToJSON(book)

	if withBalance {
		balanceMap, _ := b.GetBalance(bookId, "", "", nil)
		result["balance"] = balanceMap
	}
	return result, nil
}

func (b *BookService) GetBalance(bookId, assetId, operationType string, tx *gorm.DB) (map[string]interface{}, error) {
	balances, err := b.BookBalanceRepository.GetBalance(bookId, assetId, operationType, tx)
	// If error, return error
	if err != nil {
		logger.Logger.Errorf("Fetching Operation Failed, error: %+v", err)
		return nil, err
	}

	// If balance is found, return operation
	if balances != nil {
		balance := map[string]interface{}{}
		// groupBy assetId, ex json: {"inr": {}, "btc": {}}
		for _, balanceStruct := range *balances {
			bMap := util.StructToJSON(balanceStruct)
			bMap["balance"] = strconv.FormatFloat(balanceStruct.Balance, 'f', -1, 64)
			delete(bMap, "id")
			delete(bMap, "createdAt")
			delete(bMap, "updatedAt")
			balance[balanceStruct.AssetId] = bMap
		}

		bInterface := util.StructToJSON(&balance)
		return bInterface, nil
	}
	// else, return (empty map, nil) if no error and no balance is found.
	return map[string]interface{}{}, nil
}
