package subpub

// UnsubscribeAll - удаляет всех подписчиков для указанного subject
func (sp *SubPub) UnsubscribeAll(subject string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	// Быстрая проверка на закрытие
	if sp.closed || sp.subscribers == nil {
		return
	}

	subscribers, exists := sp.subscribers[subject]
	if !exists || subscribers == nil {
		return
	}

	// Создаем копию для безопасной итерации
	subsCopy := make([]*Subscriber, 0, len(subscribers))
	for sub := range subscribers {
		subsCopy = append(subsCopy, sub)
	}

	// Отписываем всех подписчиков
	for _, sub := range subsCopy {
		sub.Unsubscribe()
	}

	// Удаляем запись о subject
	delete(sp.subscribers, subject)
}
