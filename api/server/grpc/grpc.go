package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	pb "general_ledger_golang/api/proto/code/go"
	"general_ledger_golang/pkg/logger"
	"general_ledger_golang/pkg/util"
)

type Grpc struct {
	pb.LegerServiceServer
}

// RegisterGrpcServer will create a grpc server on the given port.
// Supports graceful stopping of the server as well.
func RegisterGrpcServer(grpcPort int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		logger.Logger.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterLegerServiceServer(s, &Grpc{})
	logger.Logger.Infof("Grpc server listening at %v", lis.Addr())

	// gracefully stopping logic...
	go util.GracefulShutDownGrpc(s)

	if err := s.Serve(lis); err != nil && err != grpc.ErrServerStopped {
		logger.Logger.Fatalf("failed to serve: %v", err)
	}
}
