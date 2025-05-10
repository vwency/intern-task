package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/intern-task/internal/service"
	"github.com/vwency/intern-task/proto/subpub"
)

type Endpoints struct {
	Subscribe endpoint.Endpoint
	Publish   endpoint.Endpoint
}

// Измените параметр на указатель
func MakeEndpoints(s *service.SubPubService) Endpoints {
	return Endpoints{
		Subscribe: makeSubscribeEndpoint(s),
		Publish:   makePublishEndpoint(s),
	}
}

func makeSubscribeEndpoint(s *service.SubPubService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*subpub.SubscribeRequest)
		msgChan, err := s.Subscribe(ctx, req.Topic)
		if err != nil {
			return nil, err
		}

		// Create a channel for gRPC stream
		stream := make(chan *subpub.Message)
		go func() {
			defer close(stream)
			for msg := range msgChan {
				stream <- &subpub.Message{
					Content: msg,
					Topic:   req.Topic,
				}
			}
		}()

		return stream, nil
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
