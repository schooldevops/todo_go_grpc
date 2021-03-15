package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/schooldevops/project/go_todo_grpc/proto/todopb"

	"google.golang.org/grpc"
)

const (
	portNumber = "9000"
	grpcServerPortNumber = "9001"
)

func main() {

	ctx := context.Background()
	mux := runtime.NewServeMux()
	options := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	if err := todopb.RegisterTodoServiceHandlerFromEndpoint(
		ctx, 
		mux,
		"localhost:" + grpcServerPortNumber,
		options,
	); err != nil {
		log.Fatalf("failed to register gRPC gateway: %v", err)
	}

	log.Printf("start HTTP on %s port", portNumber)
	if err := http.ListenAndServe(":" + portNumber, mux); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}