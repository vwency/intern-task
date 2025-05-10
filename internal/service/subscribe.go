package service

import (
	"context"
)

func (s *SubPubService) Subscribe(ctx context.Context, topic string) (chan string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil, context.Canceled
	}

	msgChan := make(chan string, 10)

	if _, exists := s.streams[topic]; !exists {
		s.streams[topic] = make(map[chan string]struct{})
	}
	s.streams[topic][msgChan] = struct{}{}

	sub := s.sp.Subscribe(topic, func(msg interface{}) {
		if str, ok := msg.(string); ok {
			select {
			case msgChan <- str:
			case <-ctx.Done():
			}
		}
	})

	go func() {
		<-ctx.Done()
		sub.Unsubscribe()

		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.streams[topic], msgChan)
		if len(s.streams[topic]) == 0 {
			delete(s.streams, topic)
		}
		close(msgChan)
	}()

	return msgChan, nil
}
