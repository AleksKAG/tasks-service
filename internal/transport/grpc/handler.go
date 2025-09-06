package grpc

import (
	"context"
	"fmt"
	"log"

	taskpb "github.com/AleksKAG/project-protos/proto/task"
	userpb "github.com/AleksKAG/project-protos/proto/user"
	"github.com/AleksKAG/tasks-service/internal/task"
)

type Handler struct {
	svc        *task.Service
	userClient userpb.UserServiceClient
	taskpb.UnimplementedTaskServiceServer
}

func NewHandler(svc *task.Service, userClient userpb.UserServiceClient) *Handler {
	return &Handler{svc: svc, userClient: userClient}
}

func (h *Handler) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	if err := ctx.Err(); err != nil {
		log.Printf("CreateTask: context error: %v", err)
		return nil, err
	}

	// Проверяем, что пользователь существует
	if _, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.UserId}); err != nil {
		log.Printf("CreateTask: user %d not found: %v", req.UserId, err)
		return nil, fmt.Errorf("user %d not found: %w", req.UserId, err)
	}

	t, err := h.svc.CreateTask(task.Task{
		UserID: req.UserId, // uint32, соответствует task.Task.UserID
		Title:  req.Title,
	})
	if err != nil {
		log.Printf("CreateTask: failed to create task for user %d: %v", req.UserId, err)
		return nil, err
	}

	return &taskpb.CreateTaskResponse{
		Task: &taskpb.Task{
			Id:     uint32(t.ID), // gorm.Model.ID (uint) -> uint32
			UserId: t.UserID,     // uint32
			Title:  t.Title,
			IsDone: t.IsDone,
		},
	}, nil
}

func (h *Handler) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.GetTaskResponse, error) {
	if err := ctx.Err(); err != nil {
		log.Printf("GetTask: context error: %v", err)
		return nil, err
	}

	log.Printf("GetTask: called with ID: %d", req.Id)
	t, err := h.svc.GetTask(uint(req.Id)) // req.Id (uint32) -> uint
	if err != nil {
		log.Printf("GetTask: error for ID %d: %v", req.Id, err)
		return nil, err
	}

	return &taskpb.GetTaskResponse{
		Task: &taskpb.Task{
			Id:     uint32(t.ID), // gorm.Model.ID (uint) -> uint32
			UserId: t.UserID,     // uint32
			Title:  t.Title,
			IsDone: t.IsDone,
		},
	}, nil
}

func (h *Handler) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.UpdateTaskResponse, error) {
	if err := ctx.Err(); err != nil {
		log.Printf("UpdateTask: context error: %v", err)
		return nil, err
	}

	// Проверяем, что пользователь существует
	if _, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.UserId}); err != nil {
		log.Printf("UpdateTask: user %d not found: %v", req.UserId, err)
		return nil, fmt.Errorf("user %d not found: %w", req.UserId, err)
	}

	// Передаём ID отдельно, так как он из gorm.Model
	t, err := h.svc.UpdateTask(uint(req.Id), task.Task{
		UserID: req.UserId, // uint32
		Title:  req.Title,
		IsDone: req.IsDone,
	})
	if err != nil {
		log.Printf("UpdateTask: failed to update task %d for user %d: %v", req.Id, req.UserId, err)
		return nil, err
	}

	return &taskpb.UpdateTaskResponse{
		Task: &taskpb.Task{
			Id:     uint32(t.ID), // gorm.Model.ID (uint) -> uint32
			UserId: t.UserID,     // uint32
			Title:  t.Title,
			IsDone: t.IsDone,
		},
	}, nil
}

func (h *Handler) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	if err := ctx.Err(); err != nil {
		log.Printf("DeleteTask: context error: %v", err)
		return nil, err
	}

	if err := h.svc.DeleteTask(uint(req.Id)); err != nil {
		log.Printf("DeleteTask: failed to delete task %d: %v", req.Id, err)
		return nil, err
	}

	return &taskpb.DeleteTaskResponse{Success: true}, nil
}

func (h *Handler) ListTasks(ctx context.Context, req *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	if err := ctx.Err(); err != nil {
		log.Printf("ListTasks: context error: %v", err)
		return nil, err
	}

	tasks, err := h.svc.ListTasks(int(req.Page), int(req.PageSize))
	if err != nil {
		log.Printf("ListTasks: failed to list tasks: %v", err)
		return nil, err
	}

	pbTasks := make([]*taskpb.Task, 0, len(tasks))
	for _, t := range tasks {
		pbTasks = append(pbTasks, &taskpb.Task{
			Id:     uint32(t.ID), // gorm.Model.ID (uint) -> uint32
			UserId: t.UserID,     // uint32
			Title:  t.Title,
			IsDone: t.IsDone,
		})
	}

	return &taskpb.ListTasksResponse{Tasks: pbTasks}, nil
}

func (h *Handler) ListTasksByUser(ctx context.Context, req *taskpb.ListTasksByUserRequest) (*taskpb.ListTasksByUserResponse, error) {
	if err := ctx.Err(); err != nil {
		log.Printf("ListTasksByUser: context error: %v", err)
		return nil, err
	}

	// Проверяем, что пользователь существует
	if _, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.UserId}); err != nil {
		log.Printf("ListTasksByUser: user %d not found: %v", req.UserId, err)
		return nil, fmt.Errorf("user %d not found: %w", req.UserId, err)
	}

	tasks, err := h.svc.ListTasksByUser(uint(req.UserId), int(req.Page), int(req.PageSize))
	if err != nil {
		log.Printf("ListTasksByUser: failed to list tasks for user %d: %v", req.UserId, err)
		return nil, err
	}

	pbTasks := make([]*taskpb.Task, 0, len(tasks))
	for _, t := range tasks {
		pbTasks = append(pbTasks, &taskpb.Task{
			Id:     uint32(t.ID), // gorm.Model.ID (uint) -> uint32
			UserId: t.UserID,     // uint32
			Title:  t.Title,
			IsDone: t.IsDone,
		})
	}

	return &taskpb.ListTasksByUserResponse{Tasks: pbTasks}, nil
}
