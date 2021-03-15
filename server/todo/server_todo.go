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
	"gopkg.in/mgo.v2/bson"
)

type todoItem struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     string             `bson:"user_id,omitempty"`
	Subject    string             `bson:"subject,omitempty"`
	Status     string             `bson:"status,omitempty"`
	CreatedAt  time.Time          `bson:"created_at,omitempty"`
	ModifiedAt time.Time          `bson:"modified_at,omitempty"`
	DoneAt     time.Time          `bson:"done_at,omitempty"`
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

	todoItem := protoTodoToItem(todo)

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

	todoItem.ID = oid
	return itemTodoToProto(todoItem), nil
}

func (*todoServer) TodoById(ctx context.Context, req *todopb.TodoCriteria) (*todopb.Todo, error) {

	id := req.GetId()

	// fmt.Printf("Read ID from req %v\n", req)
	// fmt.Printf("Read ID from %v\n", id)
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	todo := &todoItem{}
	filter := bson.M{"_id": oid}

	res := todoCollection.FindOne(context.Background(), filter)
	if err := res.Decode(todo); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find todo with specified ID: %v", err),
		)
	}

	// fmt.Printf("Read Value res: %v\n", res)
	// fmt.Printf("Read Value: %v\n", todo)
	return itemTodoToProto(*todo), nil
}

func (*todoServer) TodoByCriteriaGrpc(req *todopb.TodoCriteria, stream todopb.TodoService_TodoByCriteriaGrpcServer) error {

	fmt.Printf("todo list by criteria %v\n", req)

	// findOptions := options.Find()
	// if req.GetTimeRange() != "" {
	// 	findOptions.SetSort(bson.D{{req.GetTimeRange(), -1}})
	// }

	filter := makeFilter(req)

	cur, err := todoCollection.Find(context.Background(), filter)
	if err != nil {
		log.Fatalf("Cannot find todo list by creteria %v", req)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		todo := &todoItem{}
		err := cur.Decode(todo)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while decoding data from MongoDB: %v", err),
			)
		}
		stream.Send(itemTodoToProto(*todo))
	}

	if err := cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}
	return nil
}

func (*todoServer) TodoByCriteria(ctx context.Context, req *todopb.TodoCriteria) (*todopb.TodoList, error) {

	fmt.Printf("todo list by criteria %v\n", req)

	// findOptions := options.Find()
	// if req.GetTimeRange() != "" {
	// 	findOptions.SetSort(bson.D{{req.GetTimeRange(), -1}})
	// }

	var results []*todopb.Todo

	filter := makeFilter(req)

	fmt.Printf("Filter : %v\n", filter)
	cur, err := todoCollection.Find(context.Background(), filter)
	if err != nil {
		log.Fatalf("Cannot find todo list by creteria %v", req)
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		todo := &todoItem{}
		err := cur.Decode(todo)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while decoding data from MongoDB: %v", err),
			)
		}
		results = append(results, itemTodoToProto(*todo))
	}

	if err := cur.Err(); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}

	fmt.Printf("TodoValues : %v\n", results)

	return &todopb.TodoList{Todo: results}, nil
}

func (*todoServer) DeleteTodoByCriteria(ctx context.Context, req *todopb.TodoId) (*todopb.TodoId, error) {

	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		fmt.Printf("Cannot convert todoId %v\n", err)
	}

	filter := bson.M{"_id": oid}

	res, err := todoCollection.DeleteOne(context.Background(), filter)

	if err != nil {
		fmt.Printf("Cannot Delete todo by id %v\n", err)
		return nil, nil
	}

	if res.DeletedCount == 0 {
		return nil, nil
	}

	return &todopb.TodoId{Id: req.GetId()}, nil
}

func (*todoServer) UpdateTodoById(ctx context.Context, req *todopb.CreateTodoRequest) (*todopb.Todo, error) {
	protoTodo := req.GetTodo()
	fmt.Printf("ID: %v\n", protoTodo)
	oid, err := primitive.ObjectIDFromHex(protoTodo.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	todo := &todoItem{}
	filter := bson.M{"_id": oid}

	res := todoCollection.FindOne(context.Background(), filter)
	if err := res.Decode(todo); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	if protoTodo.GetUserId() != "" {
		todo.UserID = protoTodo.GetUserId()
	}

	if protoTodo.GetSubject() != "" {
		todo.Subject = protoTodo.GetSubject()
	}

	if protoTodo.GetStatus() != "" {
		todo.Status = protoTodo.GetStatus()
		if protoTodo.GetStatus() == "DONE" {
			todo.DoneAt = time.Now()
		}
	}

	// if protoTodo.GetModifiedAt() != nil {
	// 	fmt.Printf("Modified At: %v\n", protoTodo.GetModifiedAt())
	// 	// todo.ModifiedAt = protoTodo.GetModifiedAt().AsTime()
	// }
	todo.ModifiedAt = time.Now()

	// if protoTodo.GetDoneAt() != nil {
	// 	todo.DoneAt = protoTodo.GetDoneAt().AsTime()
	// }

	_, updateErr := todoCollection.ReplaceOne(context.Background(), filter, todo)

	fmt.Printf("Value : %v, %v\n", filter, todo)

	if updateErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot update object in MongoDB: %v", updateErr),
		)
	}

	return itemTodoToProto(*todo), nil

}
