package repository

import (
	"context"

	"gorm.io/gorm"

	"general_ledger_golang/domain"
	"general_ledger_golang/models"
)

type bookBalanceGormRepository struct {
	db *gorm.DB
}

func NewBookBalanceRepository(db *gorm.DB) BookBalanceRepository {
	return &bookBalanceGormRepository{db: db}
}

func (r *bookBalanceGormRepository) ModifyBalance(ctx context.Context, tx *gorm.DB, entries []domain.OperationEntry, metadata map[string]interface{}) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	balance := models.BookBalance{}
	return balance.ModifyBalance(ctx, entries, metadata, db)
}

func (r *bookBalanceGormRepository) GetBalance(ctx context.Context, tx *gorm.DB, bookId, assetId, operationType string) ([]models.BookBalance, error) {
	db := r.db
	if tx != nil {
		db = tx
	}
	balance := models.BookBalance{}
	return balance.GetBalance(ctx, bookId, assetId, operationType, db)
}
