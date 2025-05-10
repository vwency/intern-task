package main

import (
	"log"
	"net"

	"github.com/vwency/intern-task/internal/subpub/endpoints"
	"github.com/vwency/intern-task/internal/subpub/service"
	transportGrpc "github.com/vwency/intern-task/internal/subpub/transport/grpc"
	"github.com/vwency/intern-task/pkg/subpub"
	grpcLib "google.golang.org/grpc"
)

func main() {
	core := subpub.NewSubPub()
	defer core.Close()

	svc := service.New(core)

	eps := endpoints.MakeEndpoints(svc)

	grpcServer := grpcLib.NewServer()
	transportGrpc.RegisterGRPCServer(grpcServer, eps)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Starting test gRPC server on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
