package todo

import (
	"time"

	"github.com/schooldevops/project/go_todo_grpc/proto/todopb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/mgo.v2/bson"
)

func protoTodoToItem(todo *todopb.Todo) todoItem {
	var createdAt time.Time
	var modifiedAt time.Time
	if todo.GetCreatedAt() != nil {
		createdAt = todo.GetCreatedAt().AsTime()
	} else {
		createdAt = time.Now()
	}

	if todo.GetModifiedAt() != nil {
		modifiedAt = todo.GetModifiedAt().AsTime()
	} else {
		modifiedAt = time.Now()
	}

	todoItem := todoItem{
		UserID:     todo.GetUserId(),
		Subject:    todo.GetSubject(),
		Status:     todo.GetStatus(),
		CreatedAt:  createdAt,
		ModifiedAt: modifiedAt,
	}

	return todoItem
}

func itemTodoToProto(item todoItem) *todopb.Todo {
	var createdAt *timestamppb.Timestamp = nil
	if !item.CreatedAt.IsZero() {
		createdAt = timestamppb.New(item.CreatedAt)
	}
	var modifiedAt *timestamppb.Timestamp = nil
	if !item.ModifiedAt.IsZero() {
		modifiedAt = timestamppb.New(item.ModifiedAt)
	}
	var doneAt *timestamppb.Timestamp = nil
	if !item.DoneAt.IsZero() {
		doneAt = timestamppb.New(item.DoneAt)
	}

	return &todopb.Todo{
		Id:         item.ID.Hex(),
		UserId:     item.UserID,
		Subject:    item.Subject,
		Status:     item.Status,
		CreatedAt:  createdAt,
		ModifiedAt: modifiedAt,
		DoneAt:     doneAt,
	}
}

func makeFilter(todo *todopb.TodoCriteria) bson.M {
	filter := bson.M{}

	if todo.GetUserId() != "" {
		filter["user_id"] = todo.GetUserId()
	}

	if todo.GetSubject() != "" {
		filter["subject"] = todo.GetSubject()
	}

	if todo.GetStatus() != "" {
		filter["status"] = todo.GetStatus()
	}

	if todo.GetTimeRange() != "" {
		timeFilter := bson.M{}
		timeFilter["$gt"] = todo.GetStartTime().AsTime()
		timeFilter["$lt"] = todo.GetEndTime().AsTime()
		filter[todo.GetTimeRange()] = &timeFilter
	}

	return filter
}
