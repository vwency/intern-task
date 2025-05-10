package subpub

func (sp *SubPub) Close() error {
	sp.mu.Lock()
	if sp.closed {
		sp.mu.Unlock()
		return nil
	}
	sp.closed = true

	subscribersCopy := make(map[string]map[*Subscriber]struct{})
	for subj, subs := range sp.subscribers {
		subscribersCopy[subj] = subs
	}
	sp.subscribers = make(map[string]map[*Subscriber]struct{})
	sp.mu.Unlock()

	sp.cancel()

	for _, subs := range subscribersCopy {
		for sub := range subs {
			sub.Unsubscribe()
		}
	}

	close(sp.msgQueue)

	sp.wg.Wait()

	return nil
}
