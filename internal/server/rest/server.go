package rest

import (
	"context"
	"github.com/felixge/httpsnoop"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"rusprofile/internal/api"
)

// Runs REST gateway proxy server to gRPC server with insecure connection.
// It also register handler for swagger UI static files with path /swagger.
func Run(grpcServiceAddr string, restServiceAddr string) {
	grpcMux := runtime.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := api.RegisterRusprofileHandlerFromEndpoint(ctx, grpcMux, grpcServiceAddr, opts)
	if err != nil {
		log.Fatalf("failed register gateway handlers: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", withLogger(grpcMux))

	fs := http.FileServer(http.Dir("./swagger-ui"))
	mux.Handle("/swagger/", http.StripPrefix("/swagger", fs))

	log.Printf("REST server is listening at %s", restServiceAddr)
	if err := http.ListenAndServe(restServiceAddr, mux); err != nil {
		log.Fatalf("REST server failed to start: %v", err)
	}
}

// Logging middleware for REST requests
func withLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		m := httpsnoop.CaptureMetrics(handler, writer, request)
		log.Printf("http[%d]-- %s -- %s\n", m.Code, m.Duration, request.URL.Path)
	})
}
