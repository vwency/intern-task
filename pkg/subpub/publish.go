package subpub

import (
	"context"
)

func (sp *SubPub) Publish(subject string, msg interface{}) error {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	if sp.closed {
		return context.Canceled
	}

	select {
	case sp.msgQueue <- messageWithSubject{subject: subject, msg: msg}:
		return nil
	case <-sp.ctx.Done():
		return sp.ctx.Err()
	}
}
