package subpub

import (
	"context"
	"sync"
)

// Subscriber - структура подписчика
type Subscriber struct {
	ch     chan interface{}   // канал для сообщений
	cancel context.CancelFunc // функция отмены контекста
	wg     sync.WaitGroup     // для ожидания завершения обработки
}

// Unsubscribe - отписывает подписчика
func (sub *Subscriber) Unsubscribe() {
	if sub == nil {
		return
	}
	sub.cancel()
	close(sub.ch)
	sub.wg.Wait()
}
