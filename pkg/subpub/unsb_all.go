package subpub

// UnsubscribeAll - удаляет всех подписчиков для указанного subject
func (sp *SubPub) UnsubscribeAll(subject string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	subscribers, exists := sp.subscribers[subject]
	if !exists {
		return
	}

	for sub := range subscribers {
		sub.cancel()  // Отменяем контекст подписчика
		close(sub.ch) // Закрываем канал
		delete(subscribers, sub)
	}

	delete(sp.subscribers, subject)
}
