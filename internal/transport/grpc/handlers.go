package grpc

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	subpubv1 "github.com/vwency/intern-task/proto/subpub"
)

func makeSubscribeStreamHandler(ep endpoint.Endpoint) func(*subpubv1.SubscribeRequest, subpubv1.SubPubService_SubscribeServer) error {
	return func(req *subpubv1.SubscribeRequest, stream subpubv1.SubPubService_SubscribeServer) error {
		ctx := stream.Context()
		resp, err := ep(ctx, req)
		if err != nil {
			return err
		}

		msgChan, ok := resp.(<-chan *subpubv1.Message)
		if !ok {
			return fmt.Errorf("invalid type returned from endpoint")
		}

		for {
			select {
			case <-ctx.Done():
				return nil
			case msg, ok := <-msgChan:
				if !ok {
					return nil
				}
				if err := stream.Send(msg); err != nil {
					return err
				}
			}
		}
	}
}

// Publish handlers
func decodePublishRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*subpubv1.PublishRequest)
	if !ok {
		return nil, fmt.Errorf("invalid PublishRequest")
	}
	return req, nil
}

func encodePublishResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*subpubv1.PublishResponse)
	if !ok {
		return nil, fmt.Errorf("invalid PublishResponse")
	}
	return resp, nil
}
