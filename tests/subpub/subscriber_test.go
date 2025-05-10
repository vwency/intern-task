package tests

import (
	"context"
	"sync"
	"testing"

	"github.com/vwency/intern-task/pkg/subpub" // Замените на правильный путь к пакету
)

func TestCloseWithSubscribers(t *testing.T) {
	sp := subpub.NewSubPub()

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

	sp.Subscribe("testSubject", handler1)
	sp.Subscribe("testSubject", handler2)

	wg.Add(2)

	err := sp.Publish("testSubject", "Message for all subscribers")
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	wg.Wait()

	if msgReceived1 != "Message for all subscribers" {
		t.Fatalf("Subscriber 1 didn't receive the message")
	}
	if msgReceived2 != "Message for all subscribers" {
		t.Fatalf("Subscriber 2 didn't receive the message")
	}

	sp.Close()

	err = sp.Publish("testSubject", "This should not be sent")
	if err == nil {
		t.Fatalf("Expected context.Canceled, but message was published after close")
	}

	if err != context.Canceled {
		t.Fatalf("Expected context.Canceled, but got %v", err)
	}
}
