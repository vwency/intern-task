package service

import (
	"context"
)

func (s *SubPubService) Subscribe(ctx context.Context, topic string) (chan string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-s.ctx.Done():
		return nil, context.Canceled
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if s.closed {
		return nil, context.Canceled
	}

	msgChan := make(chan string, 10)

	if _, exists := s.streams[topic]; !exists {
		s.streams[topic] = make(map[chan string]struct{})
	}
	s.streams[topic][msgChan] = struct{}{}

	sub := s.sp.Subscribe(topic, func(msg interface{}) {
		select {
		case msgChan <- msg.(string):
		case <-ctx.Done():
		case <-s.ctx.Done():
		}
	})

	go func() {
		select {
		case <-ctx.Done():
			s.mu.Lock()
			defer s.mu.Unlock()

			if _, exists := s.streams[topic]; exists {
				if _, chExists := s.streams[topic][msgChan]; chExists {
					delete(s.streams[topic], msgChan)
					if len(s.streams[topic]) == 0 {
						delete(s.streams, topic)
					}
					close(msgChan)
				}
			}
			sub.Unsubscribe()

		case <-s.ctx.Done():
		}
	}()

	return msgChan, nil
}
