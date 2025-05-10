package grpc

import (
	"context"

	subpubv1 "github.com/vwency/intern-task/proto/subpub"
)

func decodeUnsubscribeRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*subpubv1.UnsubscribeRequest)
	if !ok {
		return nil, ErrInvalidRequestType
	}
	return req, nil
}

func encodeUnsubscribeResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(*subpubv1.UnsubscribeResponse)
	if !ok {
		return nil, ErrInvalidResponseType
	}
	return resp, nil
}

func (s *grpcServer) Unsubscribe(ctx context.Context, req *subpubv1.UnsubscribeRequest) (*subpubv1.UnsubscribeResponse, error) {
	_, resp, err := s.unsubscribe.ServeGRPC(ctx, req)
	if err != nil {
		return nil, convertToGRPCError(err)
	}
	return resp.(*subpubv1.UnsubscribeResponse), nil
}
