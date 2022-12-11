package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "general_ledger_golang/api/proto/code/go"
	"general_ledger_golang/pkg/logger"
)

func AddHeaderInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		send, _ := metadata.FromOutgoingContext(ctx)
		newMD := metadata.Pairs("authorization", "aDummyToken")
		ctx = metadata.NewOutgoingContext(ctx, metadata.Join(send, newMD))

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func LogInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		logger.Logger.Infof("%s was invoked with %v\n", method, req)

		headers, ok := metadata.FromOutgoingContext(ctx)

		if ok {
			logger.Logger.Infof("Sending headers: %v\n", headers)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func GetBalanceCall(c pb.LegerServiceClient) {
	logger.Logger.Infof("GetBalanceCall was invoked")
	r, err := c.GetBalance(context.Background(), &pb.GetBalanceReq{BookId: "4"})

	if err != nil {
		logger.Logger.Fatalf("Could not call get balance: %v\n", err)
	}

	logger.Logger.Infof("GetBalanceResp: %s\n", r.Balances)
}

func GetBookCall(c pb.LegerServiceClient) {
	logger.Logger.Infof("GetBalanceCall was invoked")
	r, err := c.GetBook(context.Background(), &pb.GetBookReq{BookId: "4"})

	if err != nil {
		logger.Logger.Fatalf("Could not call get balance: %v\n", err)
	}

	logger.Logger.Infof("GetBalanceResp: %s\n", r.Book)
}

func GetOperationByMemoCall(c pb.LegerServiceClient) {
	logger.Logger.Infof("GetBalanceCall was invoked")
	r, err := c.GetOperationByMemo(context.Background(), &pb.GetOperationByMemoReq{Memo: "MEMO_3"})

	if err != nil {
		logger.Logger.Fatalf("Could not call get balance: %v\n", err)
	}

	logger.Logger.Infof("GetBalanceResp: %s\n", r.GetOperation())
}

func CreateOperationCall(c pb.LegerServiceClient) {
	logger.Logger.Infof("GetBalanceCall was invoked")

	operation := &pb.CreateOperationReq{
		Type: "TRANSFER",
		Memo: "MEMO_HHG2499421",
		Entries: []*pb.Entries{
			&pb.Entries{
				Value:   "-100",
				BookId:  "1",
				AssetId: "inr",
			},
			&pb.Entries{
				Value:   "100",
				BookId:  "4",
				AssetId: "inr",
			},
		},
		Metadata: map[string]string{
			"operation": "DEPOSIT",
			"order_id":  "81233318822391995",
		},
	}
	r, err := c.CreateOperation(context.Background(), operation)

	if err != nil {
		logger.Logger.Fatalf("Could not call create operation: %v\n", err)
	}

	logger.Logger.Infof("CreateOperationCallResp: %s\n", r.GetOperation())
}
