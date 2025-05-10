package tests

import (
	"sync"
	"testing"

	"github.com/vwency/intern-task/pkg/subpub"
)

func TestUnsubscribeAll(t *testing.T) {
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

	sp.UnsubscribeAll("testSubject")

	msgReceived1, msgReceived2 = nil, nil
	err = sp.Publish("testSubject", "This message should not be received")
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	if msgReceived1 != nil || msgReceived2 != nil {
		t.Fatalf("Subscribers should not have received the message after unsubscribe")
	}
}
