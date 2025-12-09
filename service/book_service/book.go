package book_service

import (
	"context"
	"errors"
	"strconv"

	"gorm.io/gorm"

	"general_ledger_golang/dto"
	"general_ledger_golang/models"
	"general_ledger_golang/repository"
)

type BalanceSnapshot struct {
	AssetId       string `json:"assetId"`
	OperationType string `json:"operationType"`
	Balance       string `json:"balance"`
}

type BookWithBalance struct {
	Book     *dto.BookDTO      `json:"book"`
	Balances []BalanceSnapshot `json:"balances,omitempty"`
}

type BookService struct {
	books    repository.BookRepository
	balances repository.BookBalanceRepository
	db       *gorm.DB
}

func NewBookService(db *gorm.DB, books repository.BookRepository, balances repository.BookBalanceRepository) *BookService {
	if db == nil {
		db, _ = models.GetDB()
	}
	if books == nil {
		books = repository.NewBookRepository(db)
	}
	if balances == nil {
		balances = repository.NewBookBalanceRepository(db)
	}
	return &BookService{books: books, balances: balances, db: db}
}

func (b *BookService) GetBook(ctx context.Context, bookId string, withBalance bool) (*BookWithBalance, error) {
	if bookId == "" {
		return nil, errors.New("bookId is empty")
	}
	book, err := b.books.GetByID(ctx, nil, bookId)
	if err != nil || book == nil {
		return bookWithNil(book), err
	}
	bookDTO, err := dto.BookToDTO(book)
	if err != nil {
		return nil, err
	}
	result := &BookWithBalance{Book: bookDTO}
	if withBalance {
		balances, err := b.balances.GetBalance(ctx, nil, bookId, "", "")
		if err != nil {
			return nil, err
		}
		result.Balances = toBalanceSnapshots(balances)
	}
	return result, nil
}

func (b *BookService) GetBooks(ctx context.Context, bookIds []string, tx *gorm.DB) ([]dto.BookDTO, error) {
	if len(bookIds) < 1 {
		return nil, errors.New("bookIds length is empty")
	}
	books, err := b.books.GetMany(ctx, tx, bookIds)
	if err != nil {
		return nil, err
	}
	result := make([]dto.BookDTO, 0, len(books))
	for _, book := range books {
		bookDTO, err := dto.BookToDTO(&book)
		if err != nil {
			return nil, err
		}
		result = append(result, *bookDTO)
	}
	return result, nil
}

func (b *BookService) GetBalance(ctx context.Context, bookId, assetId, operationType string, tx *gorm.DB) ([]BalanceSnapshot, error) {
	balances, err := b.balances.GetBalance(ctx, tx, bookId, assetId, operationType)
	if err != nil {
		return nil, err
	}
	return toBalanceSnapshots(balances), nil
}

func (b *BookService) UpsertBook(ctx context.Context, payload dto.BookPayload) (*dto.BookDTO, string, error) {
	if err := payload.Validate(); err != nil {
		return nil, "", err
	}
	model, err := dto.BookPayloadToModel(payload)
	if err != nil {
		return nil, "", err
	}
	book, op, err := b.books.Upsert(ctx, nil, model)
	if err != nil {
		return nil, "", err
	}
	bookDTO, err := dto.BookToDTO(book)
	if err != nil {
		return nil, "", err
	}
	return bookDTO, op, nil
}

func (b *BookService) CheckBookExists(ctx context.Context, nUniqBookIds []string, tx *gorm.DB) (bool, error) {
	seen := map[string]struct{}{}
	unique := make([]string, 0, len(nUniqBookIds))
	for _, id := range nUniqBookIds {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		unique = append(unique, id)
	}
	books, err := b.GetBooks(ctx, unique, tx)
	if err != nil {
		return false, err
	}
	if len(books) < 1 {
		return false, errors.New("books with provided ids don't exist")
	}
	if len(books) != len(unique) {
		return false, errors.New("some bookIds couldn't be found")
	}
	return true, nil
}

func toBalanceSnapshots(balances []models.BookBalance) []BalanceSnapshot {
	result := make([]BalanceSnapshot, 0, len(balances))
	for _, balance := range balances {
		result = append(result, BalanceSnapshot{
			AssetId:       balance.AssetId,
			OperationType: balance.OperationType,
			Balance:       strconv.FormatFloat(balance.Balance, 'f', -1, 64),
		})
	}
	return result
}

func bookWithNil(book *models.Book) *BookWithBalance {
	if book == nil {
		return nil
	}
	bookDTO, err := dto.BookToDTO(book)
	if err != nil {
		return nil
	}
	return &BookWithBalance{Book: bookDTO}
}
