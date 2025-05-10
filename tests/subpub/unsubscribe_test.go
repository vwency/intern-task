package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/vwency/intern-task/pkg/subpub"
)

func TestUnsubscribe(t *testing.T) {
	sp := subpub.NewSubPub()
	defer sp.Close()

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

	sub1 := sp.Subscribe("testSubject", handler1)
	sub2 := sp.Subscribe("testSubject", handler2)

	wg.Add(2)
	err := sp.Publish("testSubject", "Message for all subscribers")
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	if waitWithTimeout(&wg, 1*time.Second) {
		t.Fatal("Timeout waiting for initial messages")
	}

	if msgReceived1 != "Message for all subscribers" {
		t.Fatalf("Subscriber 1 didn't receive the message")
	}
	if msgReceived2 != "Message for all subscribers" {
		t.Fatalf("Subscriber 2 didn't receive the message")
	}

	sub1.Unsubscribe()
	sub2.Unsubscribe()

	msgReceived1, msgReceived2 = nil, nil

	err = sp.Publish("testSubject", "Message after unsubscribe")
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if msgReceived1 != nil || msgReceived2 != nil {
		t.Fatalf("No subscriber should have received the message after unsubscribe")
	}
}

func waitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return false
	case <-time.After(timeout):
		return true
	}
}
