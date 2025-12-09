package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"general_ledger_golang/pkg/util"
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

func ValidatePostOperation(data map[string]interface{}) {
	var myMapSlice []map[string]interface{}
	entryRule := map[string]interface{}{
		"bookId":  "required,min=1",
		"assetId": "required,min=1",
		"value":   "required,min=1",
	}
	myMapSlice = append(myMapSlice, entryRule)

	rules := map[string]interface{}{
		"type": "required,min=3,max=20",
		"memo": "required,min=3",
		//"entries":  entryRule,
		"metadata": map[string]interface{}{},
	}

	v := validator.New()
	errs := v.ValidateMap(data, rules)
	resultErr := map[string]interface{}{}

	for k, err := range errs {
		resultErr[k] = err
		fmt.Println("Error is: ", err)
	}

	if len(errs) > 0 {
		data["valid"] = false
		data["errors"] = resultErr
		return
	}

	entries := data["entries"]

	if reflect.TypeOf(entries).Kind() == reflect.Slice {
		if reflect.TypeOf(entries).Elem().Kind() == reflect.Interface {
			for _, entry := range entries.([]interface{}) {
				e := v.ValidateMap(entry.(map[string]interface{}), entryRule)
				util.Copy(errs, e, false)
			}
		}
	}

	if len(errs) < 1 {
		return
	}
	data["valid"] = false
	data["errors"] = errs

	bs, _ := json.Marshal(errs)
	fmt.Println("Final Errors String: ", string(bs))

	return
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

func (o *Operation) CreateOperation(op map[string]interface{}, tx *gorm.DB) (*Operation, error) {
	// operations are idempotent
	var d *gorm.DB

	if tx != nil {
		d = tx
	} else {
		d = db
	}

	operation := Operation{}

	if op["entries"] != nil {
		entriesBytes, _ := json.Marshal(op["entries"])
		op["entries"] = datatypes.JSON(entriesBytes)
	}

	if op["metadata"] != nil {
		metadataBytes, _ := json.Marshal(op["metadata"])
		op["metadata"] = datatypes.JSON(metadataBytes)
	}

	r := d.Model(&o).FirstOrCreate(&operation, op)
	if r.Error != nil {
		return nil, r.Error
	}
	return &operation, nil
}

func (o *Operation) UpdateOperation(op map[string]interface{}, tx *gorm.DB) error {
	// operations are idempotent
	var d *gorm.DB

	if tx != nil {
		d = tx
	} else {
		d = db
	}

	r := d.Model(&o).Where("memo = ?", op["memo"]).Updates(&op)
	if r.Error != nil {
		return r.Error
	}
	return nil
}
