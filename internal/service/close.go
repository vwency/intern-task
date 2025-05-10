package service

func (s *SubPubService) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	s.closed = true

	if s.cancel != nil {
		s.cancel()
	}

	for topic := range s.streams {
		s.sp.UnsubscribeAll(topic)
	}

	for topic, streams := range s.streams {
		for ch := range streams {
			close(ch)
		}
		delete(s.streams, topic)
	}

	s.sp.WaitForCompletion()
}
