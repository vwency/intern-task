package service

func (s *SubPubService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrServiceClosed
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
			select {
			case _, ok := <-ch:
				if !ok {
					continue
				}
			default:
				close(ch)
			}
		}
		delete(s.streams, topic)
	}

	s.sp.WaitForCompletion()

	return nil
}
