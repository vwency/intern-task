package subpub

func (sp *SubPub) processMessages() {
	defer sp.wg.Done()
	for {
		select {
		case <-sp.ctx.Done():
			return
		case msgWithSubject, ok := <-sp.msgQueue:
			if !ok {
				return
			}

			sp.mu.RLock()
			subsForSubject, exists := sp.subscribers[msgWithSubject.subject]
			if !exists {
				sp.mu.RUnlock()
				continue
			}

			subsCopy := make([]*Subscriber, 0, len(subsForSubject))
			for sub := range subsForSubject {
				select {
				case <-sub.ctx.Done():
					continue
				default:
					subsCopy = append(subsCopy, sub)
				}
			}
			sp.mu.RUnlock()

			for _, sub := range subsCopy {
				select {
				case sub.ch <- msgWithSubject.msg:
				case <-sub.ctx.Done():
				case <-sp.ctx.Done():
					return
				}
			}
		}
	}
}
