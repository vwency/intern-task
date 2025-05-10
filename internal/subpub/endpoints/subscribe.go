package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/intern-task/internal/subpub/service"
	"github.com/vwency/intern-task/proto/subpub"
)

func makeSubscribeEndpoint(s *service.SubPubService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*subpub.SubscribeRequest)
		if !ok {
			return nil, ErrInvalidRequest
		}

		msgChan, err := s.Subscribe(ctx, req.Topic)
		if err != nil {
			return nil, ConvertServiceError(err)
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
