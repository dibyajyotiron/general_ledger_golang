package repository

import (
	"context"

	"gorm.io/gorm"

	"general_ledger_golang/models"
)

type OperationRepository interface {
	GetByMemo(ctx context.Context, tx *gorm.DB, memo string) (*models.Operation, error)
	Create(ctx context.Context, tx *gorm.DB, op *models.Operation) (*models.Operation, error)
	UpdateStatus(ctx context.Context, tx *gorm.DB, memo string, status models.Status, rejectionReason string) error
}
