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

// NewSubPub - создает новый Publisher
func NewSubPub() *SubPub {
	ctx, cancel := context.WithCancel(context.Background())
	sp := &SubPub{
		subscribers: make(map[string]map[*Subscriber]struct{}), // правильная инициализация
		msgQueue:    make(chan messageWithSubject, 100),
		ctx:         ctx,
		cancel:      cancel,
	}
	sp.wg.Add(1)
	go sp.processMessages()
	return sp
}

// Publish - публикует сообщение для определенного subject
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

			// Отправляем сообщение всем подписчикам subject
			for sub := range subsForSubject {
				select {
				case sub.ch <- msgWithSubject.msg:
				case <-sp.ctx.Done():
					sp.mu.RUnlock()
					return
				default:
					// Пропускаем медленных подписчиков
				}
			}
			sp.mu.RUnlock()
		}
	}
}
