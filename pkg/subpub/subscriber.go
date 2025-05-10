package subpub

import (
	"context"
	"sync"
	"sync/atomic"
)

type Subscriber struct {
	ch     chan interface{}
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	closed atomic.Bool
}

func (sub *Subscriber) Unsubscribe() {
	if sub == nil {
		return
	}

	// Проверяем, был ли уже закрыт подписчик
	if sub.closed.Swap(true) {
		return
	}

	sub.cancel() // Отменяем контекст

	// Безопасное закрытие канала
	if sub.ch != nil {
		close(sub.ch)
	}

	sub.wg.Wait() // Ждем завершения обработчика
}
