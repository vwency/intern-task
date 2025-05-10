package service

import (
	"context"
)

func (s *SubPubService) Unsubscribe(ctx context.Context, topic string, ch chan string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return context.Canceled
	}

	if streams, exists := s.streams[topic]; exists {
		if _, chExists := streams[ch]; chExists {
			delete(streams, ch)
			if len(streams) == 0 {
				delete(s.streams, topic)
			}

			select {
			case _, ok := <-ch:
				if !ok {
					return nil
				}
			default:
				close(ch)
			}
		}
	}

	return nil
}
