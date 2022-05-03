package main

import (
	"log"
	"net"

	"github.com/schooldevops/project/go_todo_grpc/server/todo"
	"google.golang.org/grpc"
)

const portNumber = "9001"

func main() {
	lis, err := net.Listen("tcp", ":"+portNumber)

	if err != nil {
		log.Fatalf("Fail to listen: %v\n", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	todo.TodoServer(s)

	log.Printf("start grpc server on %s port", portNumber)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}

}
