package tests

import (
	"testing"

	"github.com/vwency/intern-task/pkg/subpub" // Замените на правильный путь к пакету
)

func TestWaitForCompletion(t *testing.T) {
	sp := subpub.NewSubPub()

	// Create a subscriber
	var msgReceived interface{}
	handler := func(msg interface{}) { msgReceived = msg }
	sp.Subscribe("testSubject", handler)

	// Publish a message
	err := sp.Publish("testSubject", "Message for subscriber")
	if err != nil {
		t.Fatalf("Failed to publish message: %v", err)
	}

	// Wait for completion
	sp.WaitForCompletion()

	// Verify the message was received
	if msgReceived != "Message for subscriber" {
		t.Fatalf("Subscriber didn't receive the message")
	}
}
