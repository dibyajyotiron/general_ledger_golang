package repository

import (
	"context"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"general_ledger_golang/domain"
	"general_ledger_golang/models"
)

type postingGormRepository struct {
	db *gorm.DB
}

func NewPostingRepository(db *gorm.DB) PostingRepository {
	return &postingGormRepository{db: db}
}

func (r *postingGormRepository) BulkCreate(ctx context.Context, tx *gorm.DB, entries []domain.OperationEntry, operationId uint64, metadata datatypes.JSON) error {
	db := r.db.WithContext(ctx)
	if tx != nil {
		db = tx.WithContext(ctx)
	}
	posting := models.Posting{}
	return posting.BulkCreatePosting(entries, db, operationId, metadata)
}
