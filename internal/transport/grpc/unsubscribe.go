package grpc

import (
	"context"
	"fmt"

	subpubv1 "github.com/vwency/intern-task/proto/subpub"
)

func decodeUnsubscribeRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*subpubv1.UnsubscribeRequest)
	if !ok {
		return nil, fmt.Errorf("invalid UnsubscribeRequest")
	}
	return req, nil
}

func encodeUnsubscribeResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*subpubv1.UnsubscribeResponse)
	if !ok {
		return nil, fmt.Errorf("invalid UnsubscribeResponse")
	}
	return resp, nil
}

func (s *grpcServer) Unsubscribe(ctx context.Context, req *subpubv1.UnsubscribeRequest) (*subpubv1.UnsubscribeResponse, error) {
	_, resp, err := s.unsubscribe.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*subpubv1.UnsubscribeResponse), nil
}
