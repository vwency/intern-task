package endpoints

import (
	"errors"

	"github.com/vwency/intern-task/internal/subpub/service"
)

var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrInvalidArgument = errors.New("invalid argument")
)

func ConvertServiceError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, service.ErrInvalidTopic),
		errors.Is(err, service.ErrTopicTooLong),
		errors.Is(err, service.ErrEmptyMessage),
		errors.Is(err, service.ErrAlreadySubscribed),
		errors.Is(err, service.ErrNotSubscribed),
		errors.Is(err, service.ErrSubscriptionClosed),
		errors.Is(err, service.ErrNoSubscribers),
		errors.Is(err, service.ErrPublishFailed):

		return err

	default:
		return errors.New("internal error")
	}
}
