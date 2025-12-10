package repository

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"general_ledger_golang/domain"
)

type PostingRepository interface {
	BulkCreate(ctx context.Context, tx *gorm.DB, entries []domain.OperationEntry, operationId uint64, metadata datatypes.JSON) error
}
