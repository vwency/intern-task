package main

import (
	"context"
	"log"
	"time"

	"github.com/vwency/intern-task/proto/subpub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := subpub.NewSubPubServiceClient(conn)

	// Test Subscribe
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Start subscription in a goroutine
	go func() {
		stream, err := c.Subscribe(ctx, &subpub.SubscribeRequest{Topic: "test-topic"})
		if err != nil {
			log.Fatalf("could not subscribe: %v", err)
		}

		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("subscription error: %v", err)
				return
			}
			log.Printf("Received message: %s (topic: %s, timestamp: %d)",
				msg.GetContent(), msg.GetTopic(), msg.GetTimestamp())
		}
	}()

	// Test Publish
	for i := 0; i < 5; i++ {
		resp, err := c.Publish(ctx, &subpub.PublishRequest{
			Topic:   "test-topic",
			Message: "Hello from test client",
		})
		if err != nil {
			log.Fatalf("could not publish: %v", err)
		}
		log.Printf("Published message, subscriber count: %d", resp.GetSubscriberCount())
		time.Sleep(2 * time.Second)
	}
}
