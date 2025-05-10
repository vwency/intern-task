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
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := subpub.NewSubPubServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	subscriptionChan := make(chan struct{})
	streamDone := make(chan struct{})
	go func() {
		defer close(streamDone)

		stream, err := c.Subscribe(ctx, &subpub.SubscribeRequest{Topic: "test-topic"})
		if err != nil {
			log.Printf("could not subscribe: %v", err)
			close(subscriptionChan)
			return
		}
		close(subscriptionChan)

		for {
			msg, err := stream.Recv()
			if err != nil {
				if ctx.Err() != nil {
					log.Println("Subscription closed by context")
				} else {
					log.Printf("subscription error: %v", err)
				}
				return
			}
			log.Printf("Received message: %s (topic: %s, timestamp: %d)",
				msg.GetContent(), msg.GetTopic(), msg.GetTimestamp())
		}
	}()

	<-subscriptionChan
	log.Println("Subscription attempt completed")

	for i := 0; i < 5; i++ {
		resp, err := c.Publish(ctx, &subpub.PublishRequest{
			Topic:   "test-topic",
			Message: "Hello from test client",
		})
		if err != nil {
			log.Printf("could not publish: %v", err)
			break
		}
		log.Printf("Published message, subscriber count: %d", resp.GetSubscriberCount())
		time.Sleep(2 * time.Second)
	}

	_, err = c.Unsubscribe(ctx, &subpub.UnsubscribeRequest{Topic: "test-topic"})
	if err != nil {
		log.Printf("could not unsubscribe: %v", err)
	}

	log.Println("Initiating graceful shutdown...")
	cancel()

	select {
	case <-streamDone:
		log.Println("Stream closed gracefully")
	case <-time.After(5 * time.Second):
		log.Println("Timeout waiting for stream to close")
	}

	if err := conn.Close(); err != nil {
		log.Printf("error closing connection: %v", err)
	}

	log.Println("Test completed successfully")
}
