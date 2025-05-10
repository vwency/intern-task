package subpub

func (sp *SubPub) UnsubscribeAll(subject string) {
	subscribers, exists := sp.subscribers[subject]
	if !exists {
		return
	}

	for sub := range subscribers {
		sub.Unsubscribe()
		delete(subscribers, sub)
	}

	delete(sp.subscribers, subject)
}
