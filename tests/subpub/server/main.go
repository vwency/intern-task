package main

import (
	"context"
	"log"
	"net"

	"github.com/vwency/intern-task/internal/endpoints"
	"github.com/vwency/intern-task/internal/service"
	transportGrpc "github.com/vwency/intern-task/internal/transport/grpc" // Renamed import
	"github.com/vwency/intern-task/pkg/subpub"
	grpcLib "google.golang.org/grpc" // Renamed import to avoid conflict
)

func main() {
	// Create core logic
	core := subpub.NewSubPub()
	defer core.Close(context.Background())

	// Service layer
	svc := service.New(core)

	// Endpoints layer
	eps := endpoints.MakeEndpoints(svc)

	// Create and run gRPC server
	grpcServer := grpcLib.NewServer()                 // Using the renamed import
	transportGrpc.RegisterGRPCServer(grpcServer, eps) // Using the renamed import

	// Listen on port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Starting test gRPC server on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
