package service

import (
	"context"
)

func (s *SubPubService) Publish(ctx context.Context, topic, message string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.closed {
		return 0, context.Canceled
	}

	err := s.sp.Publish(topic, message)
	if err != nil {
		return 0, err
	}

	return len(s.streams[topic]), nil
}
