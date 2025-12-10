package repository

import (
	"context"

	"gorm.io/gorm"

	"general_ledger_golang/models"
)

type BookRepository interface {
	Upsert(ctx context.Context, tx *gorm.DB, book *models.Book) (*models.Book, string, error)
	GetByID(ctx context.Context, tx *gorm.DB, bookId string) (*models.Book, error)
	GetMany(ctx context.Context, tx *gorm.DB, bookIds []string) ([]models.Book, error)
}
