package grpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"

	proto "general_ledger_golang/api/proto/code/go"
	"general_ledger_golang/models"
	"general_ledger_golang/pkg/e"
	"general_ledger_golang/pkg/logger"
	"general_ledger_golang/pkg/util"
	"general_ledger_golang/service/operation_service"

	"general_ledger_golang/service/book_service"
)

func (*Grpc) GetBook(_ context.Context, req *proto.GetBookReq) (res *proto.GetBookRes, err error) {
	bookService := book_service.BookService{}

	result, err := bookService.GetBook(req.BookId, false)
	if err != nil {
		logger.Logger.Infof("Error Occured while calling bookService.GetBook, req: %+v, err: %+v", req, err)
		return nil, err
	}

	marshal, _ := json.Marshal(result)
	logger.Logger.Infof("Result: %+v", string(marshal))

	d, err := util.InterfaceToMapOfString(result["metadata"])
	if err != nil {
		logger.Logger.Infof("Error Occured while calling InterfaceToMapOfString, req: %+v, err: %+v", result["metadata"], err)
		return nil, err
	}
	mappedBook := &proto.BookResp{
		CreatedAt: result["createdAt"].(string),
		Id:        decimal.NewFromFloat(result["id"].(float64)).String(),
		Metadata:  d,
		Name:      result["name"].(string),
		UpdatedAt: result["updatedAt"].(string),
	}

	return &proto.GetBookRes{
		Book: mappedBook,
	}, nil
}

// GetBalance returns a map where key is the asset name, and value is the amount of that asset in that book
func (*Grpc) GetBalance(_ context.Context, req *proto.GetBalanceReq) (res *proto.GetBalanceRes, err error) {
	logger.Logger.Infof("Invoked GetBalance")
	bookService := book_service.BookService{}

	result, err := bookService.GetBalance(req.BookId, "", "", nil)
	marshal, _ := json.Marshal(result)

	logger.Logger.Infof("Result: %+v", string(marshal))

	if err != nil {
		return nil, err
	}

	resBalances := map[string]string{}
	for k, v := range result {
		vMap, er := util.InterfaceToMapOfString(v)
		// in case of err, return
		if er != nil {
			logger.Logger.Infof("Error Occured while converting interface to map of string, v: %+v, err: %+v", v, er)
			return nil, er
			//return nil, errors.New("something went wrong, we're checking")
		}
		resBalances[k] = vMap["balance"]
	}
	return &proto.GetBalanceRes{
		Balances: resBalances,
	}, nil
}

func (*Grpc) CreateOrUpdateBook(_ context.Context, req *proto.CreateUpdateBookReq) (*proto.CreateUpdateBookRes, error) {
	metadataBytes, _ := json.Marshal(req.Metadata)
	if req.Name == "" {
		return nil, e.GrpcFieldNotFound("name is required.")
	}
	book := models.Book{Name: req.Name, Metadata: datatypes.JSON(metadataBytes)}
	result, operationMessage := book.CreateOrUpdateBook(&book)
	err := result.Error
	if err != nil {
		logger.Logger.Errorf("Book creation failed: %+v", err)
		return nil, e.GrpcInternalError(
			"book.CreateOrUpdateBook",
			err,
			map[string]string{
				"name": book.Name,
				"err":  "Book creation failed:",
			},
		)
	}
	return &proto.CreateUpdateBookRes{
		Message: fmt.Sprintf("book %s successful", operationMessage),
	}, nil
}

func (*Grpc) GetOperationByMemo(_ context.Context, req *proto.GetOperationByMemoReq) (res *proto.GetOperationByMemoRes, err error) {
	opService := &operation_service.OperationService{}
	if req.Memo == "" {
		return nil, e.GrpcFieldNotFound("memo is required.")
	}
	foundOp, err := opService.GetOperation(req.Memo, nil)
	if err != nil {
		logger.Logger.Errorf("Fetching operation failed, memo: %+v, err: %+v", req.Memo, err)
		return nil, e.GrpcInternalError("opService.GetOperation", err, nil)
	}
	if foundOp == nil {
		errMsg := fmt.Sprintf("Operation with memo %s is not found", req.Memo)
		return nil, e.GrpcRecordNotFound(errMsg, "GetOperationByMemo", nil)
	}

	protoEntries, err2 := opService.EntryInterfaceToProtoEntries(foundOp["entries"])
	if err2 != nil {
		logger.Logger.Errorf("converting to proto entries failed, op: %+v, err: %+v", foundOp, err2)
		protoEntries = nil // it's nil in case of error so response can be sent.
		return nil, e.GrpcInternalError("opService.GetOperation", err2, nil)
	}

	metadata, err3 := util.InterfaceToMapOfString(foundOp["metadata"])
	if err3 != nil {
		logger.Logger.Errorf("converting metadata to interface failed, op: %+v, err: %+v", foundOp, err2)
	}

	operation := &proto.Operation{
		Memo:            foundOp["memo"].(string),
		Id:              decimal.NewFromFloat(foundOp["id"].(float64)).IntPart(),
		CreatedAt:       foundOp["createdAt"].(string),
		UpdatedAt:       foundOp["updatedAt"].(string),
		Type:            foundOp["type"].(string),
		Entries:         protoEntries,
		Status:          foundOp["status"].(string),
		RejectionReason: foundOp["rejectionReason"].(string),
		Metadata:        metadata,
		// note, postman, for some reason, doesn't show
		// metadata (empty object in pm), but it's shown
		// if made request from a raw cli based grpc client.
	}

	logger.Logger.Infof("operation: %+v", operation)

	return &proto.GetOperationByMemoRes{
		Operation: operation,
	}, nil

}

func (*Grpc) CreateOperation(_ context.Context, req *proto.CreateOperationReq) (res *proto.CreateOperationRes, err error) {
	opService := &operation_service.OperationService{}

	reqEntries := opService.ProtoEntriesToEntryInterface(req.Entries)
	metadataInterface := map[string]interface{}{}

	for key, value := range req.Metadata {
		metadataInterface[key] = value
	}
	opMap := map[string]interface{}{
		"type":     req.Type,
		"memo":     req.Memo,
		"entries":  reqEntries,
		"metadata": metadataInterface,
	}

	foundOp, err := opService.PostOperation(opMap)
	if err != nil || foundOp == nil {
		logger.Logger.Errorf("Creating Operation Failed, error: %+v", err)
		return nil, e.GrpcInternalError("Creating operation resulted in error!",
			err, nil)
	}

	protoEntries, err2 := opService.EntryInterfaceToProtoEntries(foundOp["entries"])
	if err2 != nil {
		logger.Logger.Errorf("converting to proto entries failed, op: %+v, err: %+v", foundOp, err2)
		protoEntries = nil // it's nil in case of error so response can be sent.
		return nil, e.GrpcInternalError("opService.GetOperation", err2, nil)
	}

	metadata, err3 := util.InterfaceToMapOfString(foundOp["metadata"])
	if err3 != nil {
		logger.Logger.Errorf("converting metadata to interface failed, op: %+v, err: %+v", foundOp, err2)
	}

	operation := &proto.Operation{
		Memo:            foundOp["memo"].(string),
		Id:              decimal.NewFromFloat(foundOp["id"].(float64)).IntPart(),
		CreatedAt:       foundOp["createdAt"].(string),
		UpdatedAt:       foundOp["updatedAt"].(string),
		Type:            foundOp["type"].(string),
		Entries:         protoEntries,
		Status:          foundOp["status"].(string),
		RejectionReason: foundOp["rejectionReason"].(string),
		Metadata:        metadata,
	}

	return &proto.CreateOperationRes{
		Operation: operation,
	}, nil
}
