package service

import (
	"context"
)

func (s *SubPubService) Publish(ctx context.Context, topic, message string) (int, error) {
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
		return 0, context.Canceled
	}

	if len(s.streams[topic]) == 0 {
		return 0, nil // Не возвращаем ошибку, если нет подписчиков
	}

	return len(s.streams[topic]), s.sp.Publish(topic, message)
}
