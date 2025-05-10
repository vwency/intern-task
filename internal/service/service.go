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
	ctx     context.Context
	cancel  context.CancelFunc
}

func New(sp *subpub.SubPub) *SubPubService {
	ctx, cancel := context.WithCancel(context.Background())
	return &SubPubService{
		sp:      sp,
		streams: make(map[string]map[chan string]struct{}),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (s *SubPubService) GetActiveSubscribersCount(topic string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.streams[topic])
}
