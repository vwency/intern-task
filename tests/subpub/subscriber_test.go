package tests

import (
	"context"
	"sync"
	"testing"

	"github.com/vwency/intern-task/pkg/subpub" // Замените на правильный путь к пакету
)

func TestCloseWithSubscribers(t *testing.T) {
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
	sp.Subscribe("testSubject", handler1)
	sp.Subscribe("testSubject", handler2)

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

	// Закрытие паблишера и очистка
	sp.Close()

	// Попытка опубликовать сообщение после закрытия
	err = sp.Publish("testSubject", "This should not be sent")
	if err == nil {
		t.Fatalf("Expected context.Canceled, but message was published after close")
	}

	// Проверка на ошибку закрытия
	if err != context.Canceled {
		t.Fatalf("Expected context.Canceled, but got %v", err)
	}
}
