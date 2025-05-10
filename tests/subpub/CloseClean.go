package tests

import (
	"context"
	"testing"

	"github.com/vwency/intern-task/pkg/subpub" // Замените на правильный путь к пакету
)

func TestCloseWithCleanup(t *testing.T) {
	sp := subpub.NewSubPub()

	// Create a subscriber
	var msgReceived interface{}
	handler := func(msg interface{}) { msgReceived = msg }
	sp.Subscribe("testSubject", handler)

	// Publish a message
	err := sp.Publish("testSubject", "Message before close")
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	// Verify the message is received
	if msgReceived != "Message before close" {
		t.Fatalf("Subscriber didn't receive the message")
	}

	// Close the publisher and clean up
	sp.Close()

	// Try publishing after closing
	err = sp.Publish("testSubject", "This should not be sent")
	if err != context.Canceled {
		t.Fatalf("Expected context.Canceled, but got %v", err)
	}
}
