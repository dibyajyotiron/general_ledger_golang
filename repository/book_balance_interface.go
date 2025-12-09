package repository

import (
	"context"

	"gorm.io/gorm"

	"general_ledger_golang/domain"
	"general_ledger_golang/models"
)

type BookBalanceRepository interface {
	ModifyBalance(ctx context.Context, tx *gorm.DB, entries []domain.OperationEntry, metadata map[string]interface{}) error
	GetBalance(ctx context.Context, tx *gorm.DB, bookId, assetId, operationType string) ([]models.BookBalance, error)
}
