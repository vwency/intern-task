package service

func (s *SubPubService) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	s.closed = true
	for topic, streams := range s.streams {
		for ch := range streams {
			close(ch)
		}
		delete(s.streams, topic)
	}
}
