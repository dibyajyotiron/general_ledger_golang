package models

import (
	"strconv"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"general_ledger_golang/domain"
)

type Posting struct {
	Model
	OperationId string         `gorm:"index;column:operationId" json:"operationId"`
	BookId      string         `gorm:"index;column:bookId" json:"bookId"`
	Value       string         `json:"status"`
	Metadata    datatypes.JSON `json:"metadata"`
	AssetId     string         `gorm:"index;column:assetId" json:"assetId"`
}

// BulkCreatePosting will create posting entries as a bulk from given slice of OperationEntry values.
func (p *Posting) BulkCreatePosting(postings []domain.OperationEntry, tx *gorm.DB, operationId uint64, metadata datatypes.JSON) error {
	// operations are idempotent
	var d *gorm.DB

	if tx != nil {
		d = tx
	} else {
		d = db
	}
	var postingsSlice []Posting

	for i := 0; i < len(postings); i++ {
		posting := postings[i]
		postingsSlice = append(postingsSlice, Posting{
			OperationId: strconv.Itoa(int(operationId)),
			BookId:      posting.BookId,
			Value:       posting.Value,
			Metadata:    metadata,
			AssetId:     posting.AssetId,
		})
	}

	r := d.Model(&p).Create(postingsSlice)
	if r.Error != nil {
		return r.Error
	}
	return nil
}
