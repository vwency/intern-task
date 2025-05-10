package subpub

// Close - завершает работу Publisher
func (sp *SubPub) Close() error {
	sp.mu.Lock()
	if sp.closed {
		sp.mu.Unlock()
		return nil
	}
	sp.closed = true
	sp.mu.Unlock()

	// Cancel the context to stop ongoing operations
	sp.cancel()

	// Close the message queue
	close(sp.msgQueue)

	// Unsubscribe all subscribers
	sp.mu.Lock()
	for subject := range sp.subscribers {
		sp.UnsubscribeAll(subject)
	}
	sp.mu.Unlock()

	// Wait for the processMessages goroutine to complete
	sp.wg.Wait()

	// Clear the subscribers map to prevent memory leaks
	sp.mu.Lock()
	sp.subscribers = make(map[string]map[*Subscriber]struct{})
	sp.mu.Unlock()

	return nil
}
