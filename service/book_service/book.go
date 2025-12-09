package book_service

import (
	"errors"
	"strconv"

	"github.com/thoas/go-funk"
	"gorm.io/gorm"

	"general_ledger_golang/models"
	"general_ledger_golang/pkg/logger"
	"general_ledger_golang/pkg/util"
)

type BookService struct {
	BookRepository        models.Book
	BookBalanceRepository models.BookBalance
}

// GetBook returns book details.
//
// Point to note: To unmarshal JSON into an interface value,
// Unmarshal stores float64, for JSON numbers in the interface value.
// So, any number inside the map that's of type interface is underlying float64.
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

func (b *BookService) GetBooks(bookIds []string, tx *gorm.DB) ([]map[string]interface{}, error) {
	if len(bookIds) < 1 {
		return nil, errors.New("BookIds length is empty")
	}
	books, err := b.BookRepository.GetBooks(bookIds, tx)
	if err != nil {
		return nil, err
	}
	if books == nil {
		return nil, nil
	}

	var result []map[string]interface{}
	//result := util.StructToJSON(books)

	for _, book := range *books {
		result = append(result, util.StructToJSON(book))
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

func (b *BookService) CheckBookExists(nUniqBookIds []string, tx *gorm.DB) (bool, error) {
	bookIds := funk.UniqString(nUniqBookIds)
	bookIdsProvided := len(bookIds)

	books, err := b.GetBooks(bookIds, tx)

	if err != nil {
		return false, err
	}

	if len(books) < 1 {
		return false, errors.New("books with provided ids don't exist")
	}

	if bookIdsProvided != len(books) {
		return false, errors.New("some bookIds couldn't be found")
	}

	return true, nil
}
