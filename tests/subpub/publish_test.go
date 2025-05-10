package tests

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/vwency/intern-task/pkg/subpub"
)

func TestPubSubLifecycle(t *testing.T) {
	t.Run("normal publish and subscribe", func(t *testing.T) {
		sp := subpub.NewSubPub()
		defer sp.Close()

		var (
			msgReceived interface{}
			wg          sync.WaitGroup
			timeout     = 1 * time.Second
		)

		wg.Add(1)
		handler := func(msg interface{}) {
			defer wg.Done()
			t.Logf("Received message: %v", msg)
			msgReceived = msg
		}

		sp.Subscribe("test", handler)
		err := sp.Publish("test", "test message")
		if err != nil {
			t.Fatalf("Publish failed: %v", err)
		}

		if waitWithTimeout(&wg, timeout) {
			t.Fatal("Timeout waiting for message")
		}

		if msgReceived != "test message" {
			t.Fatalf("Expected 'test message', got %v", msgReceived)
		}
	})

	t.Run("publish after close", func(t *testing.T) {
		sp := subpub.NewSubPub()

		var wg sync.WaitGroup
		wg.Add(1)
		sp.Subscribe("test", func(msg interface{}) {
			defer wg.Done()
			t.Logf("Received initial message: %v", msg)
		})

		err := sp.Publish("test", "initial message")
		if err != nil {
			t.Fatalf("Initial publish failed: %v", err)
		}

		if waitWithTimeout(&wg, 1*time.Second) {
			t.Fatal("Timeout waiting for initial message")
		}

		sp.Close()

		err = sp.Publish("test", "should fail")
		if err == nil {
			t.Fatal("Expected error when publishing after close, got nil")
		}
		if err != context.Canceled {
			t.Fatalf("Expected context.Canceled error, got %v", err)
		}
	})

	t.Run("subscribe after close", func(t *testing.T) {
		sp := subpub.NewSubPub()
		sp.Close()

		sub := sp.Subscribe("test", func(msg interface{}) {
			t.Error("Handler should not be called after close")
		})

		if sub != nil {
			t.Fatal("Expected nil subscriber when subscribing after close")
		}
	})
}
