package subpub

import "context"

// Subscribe - добавляет нового подписчика
func (sp *SubPub) Subscribe(subject string, cb MessageHandler) *Subscriber {
	subCtx, cancel := context.WithCancel(sp.ctx)
	sub := &Subscriber{
		ch:     make(chan interface{}, 10),
		ctx:    subCtx,
		cancel: cancel,
	}

	sp.mu.Lock()
	defer sp.mu.Unlock()

	if sp.closed {
		cancel()
		return nil
	}

	// Инициализируем map для subject, если его еще нет
	if _, exists := sp.subscribers[subject]; !exists {
		sp.subscribers[subject] = make(map[*Subscriber]struct{})
	}

	sp.subscribers[subject][sub] = struct{}{}

	sub.wg.Add(1)
	go func() {
		defer sub.wg.Done()
		for {
			select {
			case <-subCtx.Done():
				return
			case msg, ok := <-sub.ch:
				if !ok {
					return
				}
				cb(msg)
			}
		}
	}()

	return sub
}
