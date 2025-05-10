package service

import (
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

func (s *SubPubService) GetActiveSubscribersCount(topic string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.streams[topic])
}
