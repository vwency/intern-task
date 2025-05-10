package service

import "errors"

var (
	ErrServiceClosed      = errors.New("service is closed")
	ErrInvalidTopic       = errors.New("invalid topic format")
	ErrTopicTooLong       = errors.New("topic length exceeds limit")
	ErrEmptyMessage       = errors.New("message cannot be empty")
	ErrNoSubscribers      = errors.New("no active subscribers")
	ErrPublishFailed      = errors.New("failed to publish message")
	ErrAlreadySubscribed  = errors.New("already subscribed")
	ErrNotSubscribed      = errors.New("not subscribed")
	ErrSubscriptionClosed = errors.New("subscription channel closed")
)
