package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vwency/intern-task/internal/subpub/service"
	"github.com/vwency/intern-task/proto/subpub"
)

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
