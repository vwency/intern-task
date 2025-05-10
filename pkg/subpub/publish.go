package subpub

import (
	"context"
)

func (sp *SubPub) Publish(subject string, msg interface{}) error {
	sp.mu.RLock() // Получаем блокировку на чтение
	defer sp.mu.RUnlock()

	// Проверяем, был ли паблишер закрыт
	if sp.closed {
		return context.Canceled
	}

	// Публикуем сообщение в очередь для всех подписчиков
	select {
	case sp.msgQueue <- messageWithSubject{subject: subject, msg: msg}: // Если очередь не закрыта
		return nil
	case <-sp.ctx.Done(): // Если контекст завершен, возвращаем ошибку
		return sp.ctx.Err()
	}
}
