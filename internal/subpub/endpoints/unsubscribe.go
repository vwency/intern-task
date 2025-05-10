package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/intern-task/internal/subpub/service"
	"github.com/vwency/intern-task/proto/subpub"
)

func makeUnsubscribeEndpoint(s *service.SubPubService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*subpub.UnsubscribeRequest)
		if !ok {
			return nil, ErrInvalidRequest
		}

		ch := make(chan string)
		err = s.Unsubscribe(ctx, req.Topic, ch)
		if err != nil {
			return nil, convertEndpointError(err)
		}

		close(ch)
		return &subpub.UnsubscribeResponse{
			Success: true,
		}, nil
	}
}
