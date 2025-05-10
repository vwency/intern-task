package subpub

func (sp *SubPub) Close() error {
	sp.mu.Lock()
	// Проверяем, если паблишер уже закрыт
	if sp.closed {
		sp.mu.Unlock()
		return nil
	}
	sp.closed = true

	// Сохраняем подписчиков и очищаем мапу
	subscribersCopy := make(map[string]map[*Subscriber]struct{})
	for subj, subs := range sp.subscribers { // Изменил subject на subj
		subscribersCopy[subj] = subs
	}
	sp.subscribers = make(map[string]map[*Subscriber]struct{})
	sp.mu.Unlock() // Разблокируем перед долгими операциями

	// Отменяем контекст, чтобы прекратить новые публикации
	sp.cancel()

	// Отписываем всех подписчиков от всех тем
	for _, subs := range subscribersCopy { // Используем _ вместо subject
		for sub := range subs {
			sub.Unsubscribe()
		}
	}

	// Закрываем очередь сообщений
	close(sp.msgQueue)

	// Ожидаем завершения всех подписчиков и горутин
	sp.wg.Wait()

	return nil
}
