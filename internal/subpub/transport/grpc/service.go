package grpc

import (
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/vwency/intern-task/internal/subpub/endpoints"
	subpubv1 "github.com/vwency/intern-task/proto/subpub"
	"google.golang.org/grpc"
)

type grpcServer struct {
	subpubv1.UnimplementedSubPubServiceServer
	subscribe   func(*subpubv1.SubscribeRequest, subpubv1.SubPubService_SubscribeServer) error
	publish     kitgrpc.Handler
	unsubscribe kitgrpc.Handler
}

func RegisterGRPCServer(s *grpc.Server, eps endpoints.Endpoints) {
	srv := &grpcServer{
		subscribe: makeSubscribeStreamHandler(eps.Subscribe),
		publish: kitgrpc.NewServer(
			eps.Publish,
			decodePublishRequest,
			encodePublishResponse,
		),
		unsubscribe: kitgrpc.NewServer(
			eps.Unsubscribe,
			decodeUnsubscribeRequest,
			encodeUnsubscribeResponse,
		),
	}
	subpubv1.RegisterSubPubServiceServer(s, srv)
}
