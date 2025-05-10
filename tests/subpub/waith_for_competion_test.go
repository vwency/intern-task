package tests

import (
	"sync"
	"testing"
	"time"

	"github.com/vwency/intern-task/pkg/subpub"
)

func TestWaitForCompletion(t *testing.T) {
	sp := subpub.NewSubPub()
	defer func() {
		err := sp.Close()
		if err != nil {
			t.Fatalf("Failed to close SubPub: %v", err)
		}
	}()

	var (
		msgReceived interface{}
		wg          sync.WaitGroup
		received    = make(chan struct{})
	)

	wg.Add(1)
	handler := func(msg interface{}) {
		defer wg.Done()
		msgReceived = msg
		close(received)
	}

	// Subscribe
	sp.Subscribe("testSubject", handler)

	// Publish a message
	err := sp.Publish("testSubject", "Message for subscriber")
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	// Wait for the message to be processed with timeout
	select {
	case <-received:
		// Message processed successfully
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for message processing")
	}

	// Wait for handler goroutine
	wg.Wait()

	// Verify the message was received
	if msgReceived != "Message for subscriber" {
		t.Fatalf("Subscriber didn't receive the correct message, got: %v", msgReceived)
	}

	// sp.Close() already waits for background goroutines
}
