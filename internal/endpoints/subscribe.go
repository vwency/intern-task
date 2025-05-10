package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/intern-task/internal/service"
	"github.com/vwency/intern-task/proto/subpub"
)

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
