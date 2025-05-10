package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidRequestType  = errors.New("invalid request type")
	ErrInvalidResponseType = errors.New("invalid response type")
)

func convertToGRPCError(err error) error {
	if st, ok := status.FromError(err); ok {
		return st.Err()
	}

	switch {
	case errors.Is(err, ErrInvalidRequestType),
		errors.Is(err, ErrInvalidResponseType):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return err
	}
}
