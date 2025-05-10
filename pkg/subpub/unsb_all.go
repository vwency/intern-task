package subpub

func (sp *SubPub) UnsubscribeAll(subject string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	subscribers, exists := sp.subscribers[subject]
	if !exists {
		return
	}

	for sub := range subscribers {
		sub.cancel()
		close(sub.ch)
		delete(subscribers, sub)
	}

	delete(sp.subscribers, subject)
}
