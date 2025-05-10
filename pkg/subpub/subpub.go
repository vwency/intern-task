package subpub

import (
	"context"
	"sync"
)

type SubPub struct {
	mu          sync.RWMutex
	subscribers map[string]map[*Subscriber]struct{}
	msgQueue    chan messageWithSubject
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	closed      bool
}

func (sp *SubPub) WaitForCompletion() {
	sp.wg.Wait()
}

func NewSubPub() *SubPub {
	ctx, cancel := context.WithCancel(context.Background())
	sp := &SubPub{
		subscribers: make(map[string]map[*Subscriber]struct{}),
		msgQueue:    make(chan messageWithSubject, 100),
		ctx:         ctx,
		cancel:      cancel,
	}
	sp.wg.Add(1)
	go sp.processMessages()
	return sp
}
