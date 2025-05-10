package grpc

import (
	"github.com/go-kit/kit/endpoint"
	subpubv1 "github.com/vwency/intern-task/proto/subpub"
)

func (s *grpcServer) Subscribe(req *subpubv1.SubscribeRequest, stream subpubv1.SubPubService_SubscribeServer) error {
	return s.subscribe(req, stream)
}

func makeSubscribeStreamHandler(ep endpoint.Endpoint) func(*subpubv1.SubscribeRequest, subpubv1.SubPubService_SubscribeServer) error {
	return func(req *subpubv1.SubscribeRequest, stream subpubv1.SubPubService_SubscribeServer) error {
		ctx := stream.Context()

		resp, err := ep(ctx, req)
		if err != nil {
			return convertToGRPCError(err)
		}

		msgChan, ok := resp.(<-chan *subpubv1.Message)
		if !ok {
			return convertToGRPCError(ErrInvalidResponseType)
		}

		for {
			select {
			case <-ctx.Done():
				return convertToGRPCError(ErrInvalidRequestType)
			case msg, ok := <-msgChan:
				if !ok {
					return nil
				}
				if err := stream.Send(msg); err != nil {
					return convertToGRPCError(err)
				}
			}
		}
	}
}
