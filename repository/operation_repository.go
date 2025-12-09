package repository

import (
	"context"

	"gorm.io/gorm"

	"general_ledger_golang/models"
)

type operationGormRepository struct {
	db *gorm.DB
}

func NewOperationRepository(db *gorm.DB) OperationRepository {
	return &operationGormRepository{db: db}
}

func (r *operationGormRepository) GetByMemo(ctx context.Context, tx *gorm.DB, memo string) (*models.Operation, error) {
	db := r.db.WithContext(ctx)
	if tx != nil {
		db = tx.WithContext(ctx)
	}
	op := models.Operation{}
	res := db.Model(&models.Operation{}).Where("memo = ?", memo).Last(&op)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, nil
	}
	return &op, nil
}

func (r *operationGormRepository) Create(ctx context.Context, tx *gorm.DB, op *models.Operation) (*models.Operation, error) {
	db := r.db.WithContext(ctx)
	if tx != nil {
		db = tx.WithContext(ctx)
	}
	created := models.Operation{}
	res := db.Model(&models.Operation{}).FirstOrCreate(&created, op)
	if res.Error != nil {
		return nil, res.Error
	}
	return &created, nil
}

func (r *operationGormRepository) UpdateStatus(ctx context.Context, tx *gorm.DB, memo string, status models.Status, rejectionReason string) error {
	db := r.db.WithContext(ctx)
	if tx != nil {
		db = tx.WithContext(ctx)
	}
	updates := map[string]interface{}{
		"status": status,
	}
	if rejectionReason != "" {
		updates["rejectionReason"] = rejectionReason
	}

	res := db.Model(&models.Operation{}).Where("memo = ?", memo).Updates(updates)
	return res.Error
}
