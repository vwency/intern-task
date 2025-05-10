package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/vwency/intern-task/internal/subpub/endpoints"
	"github.com/vwency/intern-task/internal/subpub/service"
	grpcTransport "github.com/vwency/intern-task/internal/subpub/transport/grpc"
	"github.com/vwency/intern-task/pkg/config"
	"github.com/vwency/intern-task/pkg/subpub"
	"google.golang.org/grpc"
)

var Cfg config.ServiceConfig

func main() {
	env := config.DetectEnv()
	config.Init(env, "subpub", &Cfg)

	core := subpub.NewSubPub()
	defer core.Close()

	svc := service.New(core)

	eps := endpoints.MakeEndpoints(svc)

	grpcServer := grpc.NewServer()
	grpcTransport.RegisterGRPCServer(grpcServer, eps)

	lis, err := net.Listen("tcp", ":"+Cfg.App.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Printf("Starting gRPC server on :%s", Cfg.App.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}
