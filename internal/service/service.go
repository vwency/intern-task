package service

import (
	"context"
	"sync"

	"github.com/vwency/intern-task/pkg/subpub"
)

type SubPubService struct {
	sp      *subpub.SubPub
	mu      sync.RWMutex
	streams map[string]map[chan string]struct{}
	closed  bool
}

func New(sp *subpub.SubPub) *SubPubService {
	return &SubPubService{
		sp:      sp,
		streams: make(map[string]map[chan string]struct{}),
	}
}

func (s *SubPubService) Subscribe(ctx context.Context, topic string) (<-chan string, error) {
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

func (s *SubPubService) GetActiveSubscribersCount(topic string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.streams[topic])
}

func (s *SubPubService) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	s.closed = true
	for topic, streams := range s.streams {
		for ch := range streams {
			close(ch)
		}
		delete(s.streams, topic)
	}
}
