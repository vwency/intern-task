package subpub

import (
	"context"
	"sync"
)

type SubPub struct {
	mu          sync.RWMutex
	subscribers map[string]map[*Subscriber]struct{} // подписчики по subject
	msgQueue    chan messageWithSubject             // канал для сообщений с subject
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	closed      bool
}

type messageWithSubject struct {
	subject string
	msg     interface{}
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

func (sp *SubPub) processMessages() {
	defer sp.wg.Done()
	for {
		select {
		case <-sp.ctx.Done():
			return
		case msgWithSubject, ok := <-sp.msgQueue:
			if !ok {
				return
			}

			sp.mu.RLock()
			subsForSubject, exists := sp.subscribers[msgWithSubject.subject]
			if !exists {
				sp.mu.RUnlock()
				continue
			}

			subsCopy := make([]*Subscriber, 0, len(subsForSubject))
			for sub := range subsForSubject {
				select {
				case <-sub.ctx.Done():
					continue
				default:
					subsCopy = append(subsCopy, sub)
				}
			}
			sp.mu.RUnlock()

			for _, sub := range subsCopy {
				select {
				case sub.ch <- msgWithSubject.msg:
				case <-sub.ctx.Done():
				case <-sp.ctx.Done():
					return
				}
			}
		}
	}
}
