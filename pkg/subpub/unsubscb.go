package subpub

func (sub *Subscriber) Unsubscribe() {
	if sub == nil {
		return
	}

	if sub.closed.Swap(true) {
		return
	}

	sub.cancel()

	if sub.ch != nil {
		close(sub.ch)
	}

	sub.wg.Wait()
}
