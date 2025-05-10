package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/intern-task/internal/service"
	"github.com/vwency/intern-task/proto/subpub"
)

type Endpoints struct {
	Subscribe   endpoint.Endpoint
	Publish     endpoint.Endpoint
	Unsubscribe endpoint.Endpoint
}

func MakeEndpoints(s *service.SubPubService) Endpoints {
	return Endpoints{
		Subscribe:   makeSubscribeEndpoint(s),
		Publish:     makePublishEndpoint(s),
		Unsubscribe: makeUnsubscribeEndpoint(s),
	}
}

func makeSubscribeEndpoint(s *service.SubPubService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*subpub.SubscribeRequest)
		msgChan, err := s.Subscribe(ctx, req.Topic)
		if err != nil {
			return nil, err
		}

		stream := make(chan *subpub.Message)
		go func() {
			defer close(stream)
			for msg := range msgChan {
				select {
				case stream <- &subpub.Message{
					Content: msg,
					Topic:   req.Topic,
				}:
				case <-ctx.Done():
					return
				}
			}
		}()

		return (<-chan *subpub.Message)(stream), nil
	}
}

func makePublishEndpoint(s *service.SubPubService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*subpub.PublishRequest)
		count, err := s.Publish(ctx, req.Topic, req.Message)
		if err != nil {
			return nil, err
		}
		return &subpub.PublishResponse{SubscriberCount: int32(count)}, nil
	}
}

func makeUnsubscribeEndpoint(s *service.SubPubService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*subpub.UnsubscribeRequest)
		// Create a new channel to pass to the Unsubscribe method
		ch := make(chan string)

		// Call the Unsubscribe method with the context, topic, and channel
		err := s.Unsubscribe(ctx, req.Topic, ch)
		if err != nil {
			return nil, err
		}

		close(ch)

		return &subpub.UnsubscribeResponse{Success: true}, nil
	}
}
