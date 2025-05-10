package subpub

// Close - завершает работу Publisher
func (sp *SubPub) Close() error {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if sp.closed {
		return nil
	}
	sp.closed = true

	// 1. Отменяем контекст - это остановит новые операции
	sp.cancel()

	// 2. Отписываем всех подписчиков (закроет их каналы)
	for subject := range sp.subscribers {
		sp.UnsubscribeAll(subject)
	}

	// 3. Закрываем очередь сообщений
	close(sp.msgQueue)

	// 4. Ждем завершения processMessages
	sp.wg.Wait()

	return nil
}
