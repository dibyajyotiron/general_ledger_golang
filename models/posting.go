package models

import (
	"strconv"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Posting struct {
	Model
	OperationId string         `gorm:"index;column:operationId" json:"operationId"`
	BookId      string         `gorm:"index;column:bookId" json:"bookId"`
	Value       string         `json:"status"`
	Metadata    datatypes.JSON `json:"metadata"`
	AssetId     string         `gorm:"index;column:assetId" json:"assetId"`
}

// BulkCreatePosting will create posting entries as a bulk from given slice of maps.
// Conds will have "operationId" first, then "metadata" added to it. Any other field inside conds is ignored.
func (p *Posting) BulkCreatePosting(postings []interface{}, tx *gorm.DB, conds ...interface{}) error {
	// operations are idempotent
	var d *gorm.DB

	if tx != nil {
		d = tx
	} else {
		d = db
	}
	var postingsSlice []Posting

	for i := 0; i < len(postings); i++ {
		posting := postings[i].(map[string]interface{})

		operationId := posting["operationId"]
		metadata := posting["metadata"]

		if conds[0] != nil {
			operationId = conds[0]
		}
		if conds[1] != nil {
			metadata = conds[1]
		}

		postingsSlice = append(postingsSlice, Posting{
			OperationId: strconv.Itoa(int(operationId.(uint64))),
			BookId:      posting["bookId"].(string),
			Value:       posting["value"].(string),
			Metadata:    metadata.(datatypes.JSON),
			AssetId:     posting["assetId"].(string),
		})
	}

	r := d.Model(&p).Create(postingsSlice)
	if r.Error != nil {
		return r.Error
	}
	return nil
}
