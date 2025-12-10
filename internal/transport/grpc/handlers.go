package grpcserver

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	structpb "google.golang.org/protobuf/types/known/structpb"

	proto "general_ledger_golang/api/proto/code/go"
	"general_ledger_golang/dto"
	"general_ledger_golang/internal/e"
	"general_ledger_golang/internal/logger"
	"general_ledger_golang/internal/service/book"
	"general_ledger_golang/internal/service/operation"
)

func (*Grpc) GetBook(ctx context.Context, req *proto.GetBookReq) (res *proto.GetBookRes, err error) {
	bookService := book.NewBookService(nil, nil, nil)

	result, err := bookService.GetBook(ctx, req.BookId, false)
	if err != nil {
		logger.Logger.Infof("Error Occured while calling bookService.GetBook, req: %+v, err: %+v", req, err)
		return nil, err
	}
	if result == nil {
		return nil, e.GrpcRecordNotFound("book not found", "GetBook", nil)
	}
	metaStruct, _ := structpb.NewStruct(result.Book.Metadata)
	mappedBook := &proto.BookResp{
		CreatedAt: result.Book.CreatedAt.Format(time.RFC3339Nano),
		Id:        fmt.Sprintf("%d", result.Book.Id),
		Metadata:  metaStruct,
		Name:      result.Book.Name,
		UpdatedAt: result.Book.UpdatedAt.Format(time.RFC3339Nano),
	}

	return &proto.GetBookRes{
		Book: mappedBook,
	}, nil
}

// GetBalance returns a map where key is the asset name, and value is the amount of that asset in that book
func (*Grpc) GetBalance(ctx context.Context, req *proto.GetBalanceReq) (res *proto.GetBalanceRes, err error) {
	logger.Logger.Infof("Invoked GetBalance")
	bookService := book.NewBookService(nil, nil, nil)

	result, err := bookService.GetBalance(ctx, req.BookId, "", "", nil)
	if err != nil {
		return nil, err
	}

	resBalances := map[string]string{}
	for _, balance := range result {
		resBalances[balance.AssetId] = balance.Balance
	}
	return &proto.GetBalanceRes{
		Balances: resBalances,
	}, nil
}

func (*Grpc) CreateOrUpdateBook(ctx context.Context, req *proto.CreateUpdateBookReq) (*proto.CreateUpdateBookRes, error) {
	if req.Name == "" {
		return nil, e.GrpcFieldNotFound("name is required.")
	}
	bookService := book.NewBookService(nil, nil, nil)

	meta := map[string]any{}
	if req.Metadata != nil {
		meta = req.Metadata.AsMap()
	}
	payload := dto.BookPayload{Name: req.Name, Metadata: meta}
	_, operationMessage, err := bookService.UpsertBook(ctx, payload)
	if err != nil {
		logger.Logger.Errorf("Book creation failed: %+v", err)
		return nil, e.GrpcInternalError(
			"bookService.UpsertBook",
			err,
			map[string]string{
				"name": payload.Name,
			},
		)
	}
	return &proto.CreateUpdateBookRes{
		Message: fmt.Sprintf("book %s successful", operationMessage),
	}, nil
}

func (*Grpc) GetOperationByMemo(ctx context.Context, req *proto.GetOperationByMemoReq) (res *proto.GetOperationByMemoRes, err error) {
	opService := operation.NewOperationService(nil, nil, nil, nil, nil)
	if req.Memo == "" {
		return nil, e.GrpcFieldNotFound("memo is required.")
	}
	foundOp, err := opService.GetOperation(ctx, req.Memo)
	if err != nil {
		logger.Logger.Errorf("Fetching operation failed, memo: %+v, err: %+v", req.Memo, err)
		return nil, e.GrpcInternalError("opService.GetOperation", err, nil)
	}
	if foundOp == nil {
		errMsg := fmt.Sprintf("Operation with memo %s is not found", req.Memo)
		return nil, e.GrpcRecordNotFound(errMsg, "GetOperationByMemo", nil)
	}

	protoEntries := opService.EntriesToProto(foundOp.Entries)
	metaStruct, _ := structpb.NewStruct(foundOp.Metadata)

	operation := &proto.Operation{
		Memo:            foundOp.Memo,
		Id:              int64(foundOp.Id),
		CreatedAt:       foundOp.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:       foundOp.UpdatedAt.Format(time.RFC3339Nano),
		Type:            foundOp.Type,
		Entries:         protoEntries,
		Status:          foundOp.Status,
		RejectionReason: foundOp.RejectionReason,
		Metadata:        metaStruct,
		// note, postman, for some reason, doesn't show
		// metadata (empty object in pm), but it's shown
		// if made request from a raw cli based grpc client.
	}

	logger.Logger.Infof("operation: %+v", operation)

	return &proto.GetOperationByMemoRes{
		Operation: operation,
	}, nil

}

func (*Grpc) CreateOperation(ctx context.Context, req *proto.CreateOperationReq) (res *proto.CreateOperationRes, err error) {
	opService := operation.NewOperationService(nil, nil, nil, nil, nil)

	meta := map[string]interface{}{}
	if req.Metadata != nil {
		meta = req.Metadata.AsMap()
	}
	opPayload := dto.OperationPayload{
		Type:     req.Type,
		Memo:     req.Memo,
		Entries:  opService.ProtoEntriesToEntries(req.Entries),
		Metadata: meta,
	}

	foundOp, err := opService.PostOperation(ctx, opPayload)
	if err != nil || foundOp == nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Errorf("Creating Operation Failed")
		return nil, e.GrpcInternalError("Creating operation resulted in error!",
			err, nil)
	}

	metaStruct, _ := structpb.NewStruct(foundOp.Metadata)
	operation := &proto.Operation{
		Memo:            foundOp.Memo,
		Id:              int64(foundOp.Id),
		CreatedAt:       foundOp.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:       foundOp.UpdatedAt.Format(time.RFC3339Nano),
		Type:            foundOp.Type,
		Entries:         opService.EntriesToProto(foundOp.Entries),
		Status:          foundOp.Status,
		RejectionReason: foundOp.RejectionReason,
		Metadata:        metaStruct,
	}

	return &proto.CreateOperationRes{
		Operation: operation,
	}, nil
}
