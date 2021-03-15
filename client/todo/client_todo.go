package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/schooldevops/project/go_todo_grpc/proto/todopb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	fmt.Println("Todo Client Luanched.")

	cc, err := grpc.Dial("localhost:9001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	
	defer cc.Close()

	c := todopb.NewTodoServiceClient(cc)

	todo := &todopb.CreateTodoRequest{
		Todo: &todopb.Todo{
			UserId: "test",
			Subject: "have a breakfast",
			Status: "TODO",
			CreatedAt: timestamppb.New(time.Now()),
		},
	}

	result, err := c.CreateTodo(context.Background(), todo)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	fmt.Printf("Result is: %v", result)
}