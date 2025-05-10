package grpc

import (
	"context"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/vwency/intern-task/internal/endpoints"
	subpubv1 "github.com/vwency/intern-task/proto/subpub"
	"google.golang.org/grpc"
)

type grpcServer struct {
	subpubv1.UnimplementedSubPubServiceServer
	subscribe func(*subpubv1.SubscribeRequest, subpubv1.SubPubService_SubscribeServer) error
	publish   kitgrpc.Handler
}

func RegisterGRPCServer(s *grpc.Server, eps endpoints.Endpoints) {
	srv := &grpcServer{
		subscribe: makeSubscribeStreamHandler(eps.Subscribe),
		publish: kitgrpc.NewServer(
			eps.Publish,
			decodePublishRequest,
			encodePublishResponse,
		),
	}
	subpubv1.RegisterSubPubServiceServer(s, srv)
}

// Server-side streaming: implemented вручную
func (s *grpcServer) Subscribe(req *subpubv1.SubscribeRequest, stream subpubv1.SubPubService_SubscribeServer) error {
	return s.subscribe(req, stream)
}

func (s *grpcServer) Publish(ctx context.Context, req *subpubv1.PublishRequest) (*subpubv1.PublishResponse, error) {
	_, resp, err := s.publish.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*subpubv1.PublishResponse), nil
}
