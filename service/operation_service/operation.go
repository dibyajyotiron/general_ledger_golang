package operation_service

import (
	"fmt"
	"reflect"
	"time"

	"gorm.io/gorm"

	"general_ledger_golang/models"
	"general_ledger_golang/pkg/util"
)

type OperationService struct {
	OperationRepository models.Operation
}

func (o *OperationService) GetOperation(memo string, tx *gorm.DB) (map[string]interface{}, error) {
	foundOp, err := o.OperationRepository.GetOperation(memo, tx)
	// If error, return error
	if err != nil {
		fmt.Printf("Fetching Operation Failed, error: %+v", err)
		return nil, err
	}
	// If operation is found, return operation
	if foundOp != nil {
		opInterface := util.StructToJSON(foundOp)
		//appGin.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{"operation": foundOp})
		return opInterface, nil
	}
	// else, return nil, nil, if no error and no operation is found.
	return nil, nil
}

func (o *OperationService) PostOperation(op map[string]interface{}) (map[string]interface{}, error) {
	// This should call ApplyOperation
	newOp, err := o.ApplyOperation(op)

	if err != nil {
		return nil, err
	}

	opInterface := util.StructToJSON(newOp)
	return opInterface, nil
}

func (o *OperationService) ApplyOperation(op map[string]interface{}) (map[string]interface{}, error) {
	db, _ := models.GetDB()

	//tx := db.Begin() // starts the transaction
	var newOp *models.Operation
	var existingOp map[string]interface{}
	var err error

	deepCopiedOp := util.DeepCopyMap(op)
	err = db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		// apply operation with retries
		op["status"] = string(models.OperationInit)
		existingOp, err = o.GetOperation(op["memo"].(string), nil)

		if err != nil {
			return err
		}

		if existingOp != nil {
			return nil
		}

		newOp, err = o.ApplyOperationWithRetries(op, tx, 0)
		if err != nil {
			return err
		}

		fmt.Printf("Posting: %v\n", reflect.TypeOf(deepCopiedOp["entries"]))
		fmt.Printf("Posting 2: %+v\n", deepCopiedOp["entries"])
		postings := &models.Posting{}

		r := postings.BulkCreatePosting(deepCopiedOp["entries"].([]interface{}), tx, newOp.Id, newOp.Metadata)
		if r != nil {
			return r
		}
		//r := tx.Model(postings).Create(op["entries"])

		//if r.Error != nil {
		//	tx.Rollback()
		//	return nil, r.Error
		//}

		// create Book balance here, if the balance goes below 0, then rollBack the trx. else proceed

		newOp.Status = string(models.OperationApplied)
		op["status"] = string(models.OperationApplied)
		newOp.UpdatedAt = time.Time{}

		err = o.OperationRepository.UpdateOperation(op, tx)
		if err != nil {
			return err
		}
		//tx.Commit() // commits the transaction
		return nil // commits the transaction
	})

	if err != nil {
		return nil, err
	}

	if existingOp != nil {
		return existingOp, nil
	}

	opInterface := util.StructToJSON(*newOp)
	return opInterface, nil
	// tx.Rollback() // rolls back the transaction
	// return false, nil
}
func (o *OperationService) UpdateOperation(op map[string]interface{}, tx *gorm.DB) (map[string]interface{}, error) {
	var db *gorm.DB

	if tx != nil {
		db = tx
	} else {
		db, _ = models.GetDB()
	}

	// update operation
	r := db.Model(&o.OperationRepository).Updates(op).Where("memo = ?", op["memo"])
	if r.Error != nil {
		return nil, r.Error
	}
	return op, nil
}

func (o *OperationService) ApplyOperationWithRetries(op map[string]interface{}, tx *gorm.DB, rC int) (*models.Operation, error) {
	newOp, err := o.OperationRepository.CreateOperation(op, tx)
	totalRetryCount, currRetryCount := rC, 0

	for currRetryCount < totalRetryCount && err != nil {
		newOp, err = o.OperationRepository.CreateOperation(op, tx)
		currRetryCount++
	}

	if err != nil {
		return nil, err
	}
	return newOp, nil
}
