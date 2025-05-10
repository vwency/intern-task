package service

import (
	"context"
)

func (s *SubPubService) Unsubscribe(ctx context.Context, topic string, ch chan string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if s.closed {
		return ErrServiceClosed
	}

	if topic == "" {
		return ErrInvalidTopic
	}

	if ch == nil {
		return ErrSubscriptionClosed
	}

	streams, exists := s.streams[topic]
	if !exists {
		return ErrNotSubscribed
	}

	if _, chExists := streams[ch]; !chExists {
		return ErrNotSubscribed
	}

	delete(streams, ch)
	if len(streams) == 0 {
		delete(s.streams, topic)
	}

	select {
	case _, ok := <-ch:
		if !ok {
			return ErrSubscriptionClosed
		}
	default:
		close(ch)
	}

	return nil
}
