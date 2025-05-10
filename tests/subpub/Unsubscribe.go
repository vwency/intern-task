package tests

import (
	"sync"
	"testing"

	"github.com/vwency/intern-task/pkg/subpub" // Замените на правильный путь к пакету
)

func TestUnsubscribe(t *testing.T) {
	sp := subpub.NewSubPub()

	// Create multiple subscribers
	var msgReceived1, msgReceived2 interface{}
	var wg sync.WaitGroup

	handler1 := func(msg interface{}) {
		defer wg.Done()
		msgReceived1 = msg
	}

	handler2 := func(msg interface{}) {
		defer wg.Done()
		msgReceived2 = msg
	}

	// Подписка на сообщения
	sub1 := sp.Subscribe("testSubject", handler1)
	sub2 := sp.Subscribe("testSubject", handler2) // Используем sub2

	// Увеличиваем счетчик ожидания для двух подписчиков
	wg.Add(2)

	// Publish a message
	err := sp.Publish("testSubject", "Message for all subscribers")
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	// Ожидаем, пока оба подписчика получат сообщение
	wg.Wait()

	// Проверка, что оба подписчика получили сообщение
	if msgReceived1 != "Message for all subscribers" {
		t.Fatalf("Subscriber 1 didn't receive the message")
	}
	if msgReceived2 != "Message for all subscribers" {
		t.Fatalf("Subscriber 2 didn't receive the message")
	}

	// Отписываем первого подписчика
	sub1.Unsubscribe()
	// Отписываем второго подписчика, чтобы использовать sub2
	sub2.Unsubscribe()

	// Обнуляем полученные сообщения для следующего теста
	msgReceived1, msgReceived2 = nil, nil

	// Publish another message after unsubscribing both subscribers
	err = sp.Publish("testSubject", "Message after unsubscribe")
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	// Ожидаем, пока второй подписчик получит сообщение
	wg.Add(1)
	wg.Wait()

	// Проверка, что ни один подписчик не получил сообщение
	if msgReceived1 != nil || msgReceived2 != nil {
		t.Fatalf("No subscriber should have received the message after unsubscribe")
	}
}
