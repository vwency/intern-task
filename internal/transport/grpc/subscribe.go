package grpc

import (
	"fmt"

	"github.com/go-kit/kit/endpoint"
	subpubv1 "github.com/vwency/intern-task/proto/subpub"
)

func (s *grpcServer) Subscribe(req *subpubv1.SubscribeRequest, stream subpubv1.SubPubService_SubscribeServer) error {
	return s.subscribe(req, stream)
}

func makeSubscribeStreamHandler(ep endpoint.Endpoint) func(*subpubv1.SubscribeRequest, subpubv1.SubPubService_SubscribeServer) error {
	return func(req *subpubv1.SubscribeRequest, stream subpubv1.SubPubService_SubscribeServer) error {
		ctx := stream.Context()

		// Call the endpoint to get the message channel
		resp, err := ep(ctx, req)
		if err != nil {
			return err
		}

		msgChan, ok := resp.(<-chan *subpubv1.Message)
		if !ok {
			return fmt.Errorf("invalid type returned from endpoint")
		}

		// Loop over messages and send them to the gRPC stream
		for {
			select {
			case <-ctx.Done(): // If the context is canceled, stop the loop
				return nil
			case msg, ok := <-msgChan:
				if !ok {
					return nil
				}
				if err := stream.Send(msg); err != nil { // Send message to the client
					return err
				}
			}
		}
	}
}
