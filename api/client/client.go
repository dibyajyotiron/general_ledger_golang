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
	r, err := c.GetBalance(context.Background(), &pb.GetBalanceRequest{BookId: "4"})

	if err != nil {
		logger.Logger.Fatalf("Could not call get balance: %v\n", err)
	}

	logger.Logger.Infof("GetBalanceResp: %s\n", r.Balances)
}
