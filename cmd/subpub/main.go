package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/vwency/intern-task/internal/endpoints"
	"github.com/vwency/intern-task/internal/service"
	grpcTransport "github.com/vwency/intern-task/internal/transport/grpc"
	"github.com/vwency/intern-task/pkg/subpub"
	"google.golang.org/grpc"
)

func main() {
	// Создаем core-логику
	core := subpub.NewSubPub()
	defer core.Close(context.Background())

	// Сервисный слой
	svc := service.New(core)

	// Endpoints слой
	eps := endpoints.MakeEndpoints(svc)

	// Создаем и запускаем gRPC-сервер
	grpcServer := grpc.NewServer()
	grpcTransport.RegisterGRPCServer(grpcServer, eps)

	// Слушаем порт
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Println("Starting gRPC server on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}
