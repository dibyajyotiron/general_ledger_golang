package models

import (
	"errors"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Status string

const (
	OperationInit     Status = "INIT"
	OperationApplied  Status = "APPLIED"
	OperationRejected Status = "REJECTED"
)

type Operation struct {
	Model
	Type            string         `gorm:"index;" json:"type"`
	Memo            string         `gorm:"index;unique" json:"memo"`
	Entries         datatypes.JSON `json:"entries"`
	Status          string         `json:"status"`
	RejectionReason string         `gorm:"index;column:rejectionReason" json:"rejectionReason"`
	Metadata        datatypes.JSON `json:"metadata"`
}

func (o *Operation) GetOperation(memo string, tx *gorm.DB) (*Operation, error) {
	if tx != nil {
		db = tx
	}
	op := Operation{}
	res := db.Model(&o).Where("memo = ?", memo).Last(&op)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if res.Error != nil {
		return nil, res.Error
	}

	return &op, nil
}

func (o *Operation) CreateOperation(op *Operation, tx *gorm.DB) (*Operation, error) {
	var d *gorm.DB

	if tx != nil {
		d = tx
	} else {
		d = db
	}

	operation := Operation{}

	r := d.Model(&o).FirstOrCreate(&operation, op)
	if r.Error != nil {
		return nil, r.Error
	}
	return &operation, nil
}

func (o *Operation) UpdateOperationStatus(memo string, status Status, rejectionReason string, tx *gorm.DB) error {
	var d *gorm.DB

	if tx != nil {
		d = tx
	} else {
		d = db
	}

	updates := map[string]interface{}{
		"status": status,
	}
	if rejectionReason != "" {
		updates["rejectionReason"] = rejectionReason
	}

	r := d.Model(&o).Where("memo = ?", memo).Updates(updates)
	if r.Error != nil {
		return r.Error
	}
	return nil
}
