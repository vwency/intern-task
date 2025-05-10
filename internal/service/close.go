package service

func (s *SubPubService) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	s.closed = true

	// Cancel the context first to stop ongoing operations
	if s.cancel != nil {
		s.cancel()
	}

	// Unsubscribe all topics from the underlying SubPub
	for topic := range s.streams {
		s.sp.UnsubscribeAll(topic)
	}

	// Close all channels
	for topic, streams := range s.streams {
		for ch := range streams {
			close(ch)
		}
		delete(s.streams, topic)
	}

	// Make sure any ongoing message processing has finished
	s.sp.WaitForCompletion()
}
