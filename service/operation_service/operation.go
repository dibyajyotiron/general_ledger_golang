package operation_service

import (
	"fmt"
	"reflect"
	"time"

	"github.com/thoas/go-funk"
	"gorm.io/gorm"

	proto "general_ledger_golang/api/proto/code/go"

	"general_ledger_golang/models"
	"general_ledger_golang/pkg/util"
	"general_ledger_golang/service/book_service"
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

// PostOperation op-> {type: '', memo:'', entries: [{bookId: '', assetId: '', value: ''}], metadata: {}}
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

	var newOp *models.Operation

	existingOp, err := o.GetOperation(op["memo"].(string), nil)

	if err != nil {
		return nil, err
	}

	if existingOp != nil {
		return existingOp, nil
	}
	bS := book_service.BookService{}

	// TODO: Maybe create rejected status in case of non_negative_balance as well, to have a better insight.
	// This memo will not be further tried, as ledger is meant to be idempotent, a new memo should be created with correct bookIds.

	deepCopiedOp := util.DeepCopyMap(op)

	err = db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		// apply operation with retries
		// return nil commits trx, return error will roll back transaction
		op["status"] = string(models.OperationInit)

		newOp, err = o.applyOperationWithRetries(op, tx, 0)
		if err != nil {
			return err
		}

		bookIds := funk.Map(deepCopiedOp["entries"], func(entry interface{}) string {
			e := entry.(map[string]interface{})
			id := e["bookId"]
			if reflect.TypeOf(id).Kind() != reflect.String {
				return ""
			}
			return id.(string)
		})
		ok, e := bS.CheckBookExists(bookIds.([]string), tx)

		if !ok {
			op["status"] = string(models.OperationRejected)
			op["rejectionReason"] = e.Error()
			_, err = o.UpdateOperation(op, tx)
			// update newOp as that will get returned to the user.
			newOp.Status = string(models.OperationRejected)
			newOp.RejectionReason = e.Error()
			if err != nil {
				return err
			}
			// return nil to create rejected operation
			return nil
		}

		postings := &models.Posting{}

		err = postings.BulkCreatePosting(deepCopiedOp["entries"].([]interface{}), tx, newOp.Id, newOp.Metadata)
		if err != nil {
			return err
		}

		// create Book balance here, if the balance goes below 0, then rollBack the trx. else proceed
		err = bS.BookBalanceRepository.ModifyBalance(deepCopiedOp, tx)
		if err != nil {
			return err
		}

		newOp.Status = string(models.OperationApplied)
		op["status"] = string(models.OperationApplied)
		newOp.UpdatedAt = time.Time{}

		err = o.OperationRepository.UpdateOperation(op, tx)
		if err != nil {
			return err
		}

		return nil // commits the transaction
	})

	if err != nil {
		return nil, err
	}

	opInterface := util.StructToJSON(*newOp)
	return opInterface, nil
}

func (o *OperationService) UpdateOperation(op map[string]interface{}, tx *gorm.DB) (map[string]interface{}, error) {
	var db *gorm.DB

	if tx != nil {
		db = tx
	} else {
		db, _ = models.GetDB()
	}

	// update operation
	r := db.Model(&o.OperationRepository).Where("memo = ?", op["memo"]).Updates(op)
	if r.Error != nil {
		return nil, r.Error
	}
	return op, nil
}

func (o *OperationService) applyOperationWithRetries(op map[string]interface{}, tx *gorm.DB, rC int) (*models.Operation, error) {
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

func (o *OperationService) EntryInterfaceToProtoEntries(entries interface{}) ([]*proto.Entries, error) {
	var protoEntries []*proto.Entries
	entriesSlice, err := util.ConvertToMapSlice(entries)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(entriesSlice); i++ {
		entryMap, err2 := util.InterfaceToMapOfString(entriesSlice[i])
		if err2 != nil {
			continue // in case of error, log and move to next
		}
		protoEntry := &proto.Entries{
			Value:   entryMap["value"],
			BookId:  entryMap["bookId"],
			AssetId: entryMap["assetId"],
		}
		protoEntries = append(protoEntries, protoEntry)
	}
	return protoEntries, nil
}

func (o *OperationService) ProtoEntriesToEntryInterface(entries []*proto.Entries) []interface{} {
	if len(entries) < 1 {
		return []interface{}{}
	}
	var resultMap []interface{}
	for i := 0; i < len(entries); i++ {
		protoEntry := entries[i]
		mapEntry := map[string]interface{}{
			"bookId":  protoEntry.BookId,
			"assetId": protoEntry.AssetId,
			"value":   protoEntry.Value,
		}
		resultMap = append(resultMap, mapEntry)
	}
	return resultMap
}
