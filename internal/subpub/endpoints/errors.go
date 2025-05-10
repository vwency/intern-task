package endpoints

import (
	"errors"

	"github.com/vwency/intern-task/internal/subpub/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrInvalidArgument = errors.New("invalid argument")
)

func convertEndpointError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, service.ErrInvalidTopic),
		errors.Is(err, service.ErrTopicTooLong),
		errors.Is(err, service.ErrEmptyMessage):
		return status.Error(codes.InvalidArgument, "invalid input: "+err.Error())

	case errors.Is(err, service.ErrAlreadySubscribed):
		return status.Error(codes.AlreadyExists, err.Error())

	case errors.Is(err, service.ErrNotSubscribed):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, service.ErrNoSubscribers),
		errors.Is(err, service.ErrPublishFailed):
		return status.Error(codes.Internal, "publish failed: "+err.Error())

	case errors.Is(err, service.ErrSubscriptionClosed):
		return status.Error(codes.Aborted, err.Error())

	default:
		return status.Error(codes.Unknown, "internal error")
	}
}
