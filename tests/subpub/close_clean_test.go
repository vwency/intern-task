package tests

import (
	"context"
	"sync"
	"testing"

	"github.com/vwency/intern-task/pkg/subpub"
)

func TestCloseWithCleanup(t *testing.T) {
	sp := subpub.NewSubPub()
	defer sp.Close()

	var (
		msgReceived interface{}
		wg          sync.WaitGroup
	)

	wg.Add(1)
	handler := func(msg interface{}) {
		defer wg.Done()
		msgReceived = msg
	}

	sp.Subscribe("testSubject", handler)

	err := sp.Publish("testSubject", "Message before close")
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	wg.Wait()

	if msgReceived != "Message before close" {
		t.Fatalf("Subscriber didn't receive the message, got: %v", msgReceived)
	}

	sp.Close()

	err = sp.Publish("testSubject", "This should not be sent")
	if err != context.Canceled {
		t.Fatalf("Expected context.Canceled, but got %v", err)
	}
}
