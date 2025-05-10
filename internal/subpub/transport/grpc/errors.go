package grpc

import (
	"context"
	"errors"

	"github.com/vwency/intern-task/internal/subpub/endpoints"
	"github.com/vwency/intern-task/internal/subpub/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidRequestType  = errors.New("invalid request type")
	ErrInvalidResponseType = errors.New("invalid response type")
)

func ConvertToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, ErrInvalidRequestType),
		errors.Is(err, ErrInvalidResponseType),
		errors.Is(err, service.ErrInvalidTopic),
		errors.Is(err, endpoints.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())

	case errors.Is(err, service.ErrTopicTooLong),
		errors.Is(err, service.ErrEmptyMessage):
		return status.Error(codes.FailedPrecondition, err.Error())

	case errors.Is(err, service.ErrAlreadySubscribed):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, service.ErrNotSubscribed):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, service.ErrSubscriptionClosed):
		return status.Error(codes.Aborted, "subscription closed: "+err.Error())

	case errors.Is(err, service.ErrServiceClosed):
		return status.Error(codes.Unavailable, "service unavailable: "+err.Error())

	case errors.Is(err, service.ErrNoSubscribers):
		return status.Error(codes.FailedPrecondition, "no subscribers: "+err.Error())

	case errors.Is(err, service.ErrPublishFailed):
		return status.Error(codes.Internal, "publish operation failed: "+err.Error())

	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request canceled")

	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request timed out")

	default:
		if st, ok := status.FromError(err); ok {
			return st.Err()
		}
		return status.Error(codes.Unknown, "internal error: "+err.Error())
	}
}
