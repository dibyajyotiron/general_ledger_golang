package e

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcFieldNotFound(message string) error {
	st := status.New(codes.InvalidArgument, message)

	br := &errdetails.BadRequest{}
	st, err := st.WithDetails(br)
	if err != nil {
		// If this errored, it will always error
		// here, so better panic so we can figure
		// out why than have this silently passing.
		panic(fmt.Sprintf("Unexpected error: %v", err))
	}
	return st.Err()
}
func GrpcInternalError(method string, err error, metadata map[string]string) error {
	st := status.New(codes.Internal, GetMsg(ERROR))

	ie := &errdetails.ErrorInfo{
		Reason:   err.Error(),
		Domain:   method,
		Metadata: metadata,
	}
	st, err = st.WithDetails(ie)
	if err != nil {
		// If this errored, it will always error
		// here, so better panic so we can figure
		// out why than have this silently passing.
		panic(fmt.Sprintf("Unexpected error: %v", err))
	}
	return st.Err()
}
func GrpcRecordNotFound(message string, method string, metadata map[string]string) error {
	st := status.New(codes.NotFound, message)

	nfr := &errdetails.ErrorInfo{
		Reason:   message,
		Domain:   method,
		Metadata: metadata,
	}
	st, err := st.WithDetails(nfr)
	if err != nil {
		// If this errored, it will always error
		// here, so better panic so we can figure
		// out why than have this silently passing.
		panic(fmt.Sprintf("Unexpected error: %v", err))
	}
	return st.Err()
}

//func FormGrpcError(code codes.Code, message string) *status.Status {
//	st := status.New(code, "invalid username")
//	v := &errdetails.BadRequest_FieldViolation{
//		Field:       "username",
//		Description: message,
//	}
//	br := &errdetails.BadRequest{}
//	br.FieldViolations = append(br.FieldViolations, v)
//	br.
//	st, err := st.WithDetails(br)
//	if err != nil {
//		// If this errored, it will always error
//		// here, so better panic so we can figure
//		// out why than have this silently passing.
//		panic(fmt.Sprintf("Unexpected error: %v", err))
//	}
//	return st
//}
