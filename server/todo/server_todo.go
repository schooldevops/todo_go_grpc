package todo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/schooldevops/project/go_todo_grpc/db/mongodb"
	"github.com/schooldevops/project/go_todo_grpc/proto/todopb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type todoItem struct {
	ID	primitive.ObjectID `bson:"_id,omitempty`
	UserID string `bson:"user_id",omitempty`
	Subject string `bson:"subject",omitempty`
	Status string `bson:"ststus",omitempty`
	CreatedAt time.Time `bson:"created_at,omitempty`
	ModifiedAt time.Time `bson:"modified_at,omitempty`
	DoneAt time.Time `bson:"done_at,omitempty`
}

type todoServer struct {
	todopb.TodoServiceServer
}

var todoCollection *mongo.Collection

func TodoServer(server *grpc.Server) {

	todoCollection = mongodb.NewMongoCollection("todos")
	
	fmt.Println("Regist Toso Service Server")
	todopb.RegisterTodoServiceServer(server, &todoServer{})
}

func (*todoServer) CreateTodo(ctx context.Context, req *todopb.CreateTodoRequest) (*todopb.Todo, error) {
	todo := req.GetTodo()
	fmt.Printf("Receive Todo for req: %v\n", req)
	fmt.Printf("Receive Todo for create: %v\n", todo)

	todoItem := todoItem{
		UserID: todo.GetUserId(),
		Subject: todo.GetSubject(),
		Status: todo.GetStatus(),
		CreatedAt: time.Now(),
	}

	res, err := todoCollection.InsertOne(context.Background(), todoItem)

	if err != nil {
		log.Fatalf("Insert Error %v", err)
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID"),
		)
	}
	
	return &todopb.Todo{
		Id: oid.Hex(),
		UserId: todoItem.UserID,
		Subject: todoItem.Subject,
		Status: todoItem.Status,
		CreatedAt: timestamppb.New(todoItem.CreatedAt),
	}, nil
}

func (*todoServer) TodoById(ctx context.Context, req *todopb.TodoCriteria) (*todopb.Todo, error) {
	return nil, nil
}

func (*todoServer) TodoByCriteria(req *todopb.TodoCriteria, stream todopb.TodoService_TodoByCriteriaServer) error {
	return nil
}