package subpub

func (sp *SubPub) Close() error {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	// Проверяем, если паблишер уже закрыт
	if sp.closed {
		return nil
	}
	sp.closed = true

	// Отменяем контекст, чтобы прекратить новые публикации
	sp.cancel()

	// Отписываем всех подписчиков от всех тем
	for subject := range sp.subscribers {
		sp.UnsubscribeAll(subject)
	}

	// Закрываем очередь сообщений
	close(sp.msgQueue)

	// Ожидаем завершения всех подписчиков и горутин
	sp.wg.Wait()

	return nil
}
