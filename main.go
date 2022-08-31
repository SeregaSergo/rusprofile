package main

import (
	"fmt"
	"os"
	"os/signal"
	"rusprofile/internal/server/grpc"
	"rusprofile/internal/server/rest"
	"syscall"
)

var (
	grpcAddr string
	restAddr string
)

func init() {
	var ok bool
	if grpcAddr, ok = os.LookupEnv("GRPC_ADDR"); !ok {
		grpcAddr = "localhost:8081"
	}
	if restAddr, ok = os.LookupEnv("REST_GATEWAY_ADDR"); !ok {
		restAddr = "localhost:8080"
	}
}

func main() {
	gracefullShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefullShutdown, syscall.SIGINT, syscall.SIGTERM)

	grpcServ := grpc.NewServer()

	go rest.Run(grpcAddr, restAddr)
	go grpc.Run(grpcServ, grpcAddr)

	<-gracefullShutdown
	fmt.Println("Shutting down the server")
	grpcServ.GracefulStop()
}
