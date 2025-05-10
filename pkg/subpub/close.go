package subpub

import "context"

// Close - завершает работу Publisher
func (sp *SubPub) Close(ctx context.Context) error {
	sp.mu.Lock()
	if sp.closed {
		sp.mu.Unlock()
		return nil
	}
	sp.closed = true
	sp.cancel()
	close(sp.msgQueue)

	// Не закрываем каналы подписчиков здесь - они уже закрываются в Unsubscribe()
	// Просто очищаем список подписчиков
	sp.subscribers = make(map[string]map[*Subscriber]struct{})
	sp.mu.Unlock()

	done := make(chan struct{})
	go func() {
		sp.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
