package subpub

import (
	"context"
	"sync"
)

type Subscriber struct {
	ch     chan interface{}
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func (sub *Subscriber) Unsubscribe() {
	if sub == nil {
		return
	}
	sub.cancel()
	close(sub.ch)
	sub.wg.Wait()
}
