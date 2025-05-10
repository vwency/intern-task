package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"

	transportgrpc "github.com/vwency/intern-task/internal/transport/grpc"
	"github.com/vwency/intern-task/pkg/subpub"
	pb "github.com/vwency/intern-task/proto/subpub"
)

func main() {
	// Create core SubPub instance
	core := subpub.New()

	// Create gRPC service with all handlers
	grpcService := transportgrpc.NewGRPCService(core)

	// Create and configure gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterSubPubServiceServer(grpcServer, &grpcServerWrapper{
		subscribeHandler: grpcService.SubscribeHandler,
		publishHandler:   grpcService.PublishHandler,
	})

	// Start listening
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Start server in goroutine
	go func() {
		log.Println("Starting gRPC server on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server gracefully...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}

// grpcServerWrapper implements pb.SubPubServiceServer
type grpcServerWrapper struct {
	pb.UnimplementedSubPubServiceServer
	subscribeHandler kitgrpc.Handler
	publishHandler   kitgrpc.Handler
}

func (s *grpcServerWrapper) Subscribe(req *pb.SubscribeRequest, stream pb.SubPubService_SubscribeServer) error {
	_, err := s.subscribeHandler.ServeGRPC(stream.Context(), req)
	return err
}

func (s *grpcServerWrapper) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
	resp, err := s.publishHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.PublishResponse), nil
}
