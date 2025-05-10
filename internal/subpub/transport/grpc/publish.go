package grpc

import (
	"context"

	subpubv1 "github.com/vwency/intern-task/proto/subpub"
)

func decodePublishRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*subpubv1.PublishRequest)
	if !ok {
		return nil, ErrInvalidRequestType
	}
	return req, nil
}

func encodePublishResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*subpubv1.PublishResponse)
	if !ok {
		return nil, ErrInvalidResponseType
	}
	return resp, nil
}

func (s *grpcServer) Publish(ctx context.Context, req *subpubv1.PublishRequest) (*subpubv1.PublishResponse, error) {
	_, resp, err := s.publish.ServeGRPC(ctx, req)
	if err != nil {
		return nil, convertToGRPCError(err)
	}
	return resp.(*subpubv1.PublishResponse), nil
}
