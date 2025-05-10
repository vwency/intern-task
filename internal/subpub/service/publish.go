package service

import (
	"context"
)

func (s *SubPubService) Publish(ctx context.Context, topic, message string) (int, error) {
	// Проверка блокировки для чтения
	s.mu.RLock()
	defer s.mu.RUnlock()

	select {
	case <-s.ctx.Done():
		return 0, context.Canceled
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	if s.closed {
		return 0, ErrServiceClosed
	}

	if topic == "" {
		return 0, ErrInvalidTopic
	}
	if len(topic) > 255 {
		return 0, ErrTopicTooLong
	}

	if message == "" {
		return 0, ErrEmptyMessage
	}
	subscribers, exists := s.streams[topic]
	if !exists || len(subscribers) == 0 {
		return 0, ErrNoSubscribers
	}
	if err := s.sp.Publish(topic, message); err != nil {
		return 0, ErrPublishFailed
	}
	return len(subscribers), nil
}
