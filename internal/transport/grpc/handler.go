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

func NewHandler(svc *task.Service, uc userpb.UserServiceClient) *Handler {
	return &Handler{svc: svc, userClient: uc}
}

func (h *Handler) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	log.Printf("CreateTask called with userId: %d, title: %s", req.UserId, req.Title)
	// Проверка существования пользователя
	if _, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.UserId}); err != nil {
		log.Printf("GetUser failed for userId: %d, error: %v", req.UserId, err)
		return nil, fmt.Errorf("user %d not found: %w", req.UserId, err)
	}
	// Создание задачи
	t, err := h.svc.CreateTask(task.Task{
		UserID: req.UserId,
		Title:  req.Title,
	})
	if err != nil {
		log.Printf("CreateTask failed: %v", err)
		return nil, err
	}
	// Ответ
	return &taskpb.CreateTaskResponse{
		Task: &taskpb.Task{
			Id:     uint32(t.ID),
			UserId: t.UserID,
			Title:  t.Title,
			IsDone: t.IsDone,
		},
	}, nil
}

func (h *Handler) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.GetTaskResponse, error) {
	log.Printf("GetTask called with id: %d", req.Id)
	t, err := h.svc.GetTask(req.Id)
	if err != nil {
		log.Printf("GetTask failed: %v", err)
		return nil, err
	}
	return &taskpb.GetTaskResponse{
		Task: &taskpb.Task{
			Id:     uint32(t.ID),
			UserId: t.UserID,
			Title:  t.Title,
			IsDone: t.IsDone,
		},
	}, nil
}

func (h *Handler) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.UpdateTaskResponse, error) {
	log.Printf("UpdateTask called with id: %d, userId: %d", req.Id, req.UserId)
	// Проверка существования пользователя
	if _, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.UserId}); err != nil {
		log.Printf("GetUser failed for userId: %d, error: %v", req.UserId, err)
		return nil, fmt.Errorf("user %d not found: %w", req.UserId, err)
	}
	// Обновление задачи
	t := task.Task{
		UserID: req.UserId,
		Title:  req.Title,
		IsDone: req.IsDone,
	}
	t.ID = uint(req.Id)
	updatedTask, err := h.svc.UpdateTask(t)
	if err != nil {
		log.Printf("UpdateTask failed: %v", err)
		return nil, err
	}
	// Ответ
	return &taskpb.UpdateTaskResponse{
		Task: &taskpb.Task{
			Id:     uint32(updatedTask.ID),
			UserId: updatedTask.UserID,
			Title:  updatedTask.Title,
			IsDone: updatedTask.IsDone,
		},
	}, nil
}

func (h *Handler) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	log.Printf("DeleteTask called with id: %d", req.Id)
	err := h.svc.DeleteTask(req.Id)
	if err != nil {
		log.Printf("DeleteTask failed: %v", err)
		return nil, err
	}
	return &taskpb.DeleteTaskResponse{Success: true}, nil
}

func (h *Handler) ListTasks(ctx context.Context, req *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	log.Printf("ListTasks called with page: %d, pageSize: %d", req.Page, req.PageSize)
	tasks, err := h.svc.ListTasks(int(req.Page), int(req.PageSize))
	if err != nil {
		log.Printf("ListTasks failed: %v", err)
		return nil, err
	}
	pbTasks := make([]*taskpb.Task, len(tasks))
	for i, t := range tasks {
		pbTasks[i] = &taskpb.Task{
			Id:     uint32(t.ID),
			UserId: t.UserID,
			Title:  t.Title,
			IsDone: t.IsDone,
		}
	}
	return &taskpb.ListTasksResponse{Tasks: pbTasks}, nil
}

func (h *Handler) ListTasksByUser(ctx context.Context, req *taskpb.ListTasksByUserRequest) (*taskpb.ListTasksByUserResponse, error) {
	log.Printf("ListTasksByUser called with userId: %d, page: %d, pageSize: %d", req.UserId, req.Page, req.PageSize)
	// Проверка существования пользователя
	if _, err := h.userClient.GetUser(ctx, &userpb.GetUserRequest{Id: req.UserId}); err != nil {
		log.Printf("GetUser failed for userId: %d, error: %v", req.UserId, err)
		return nil, fmt.Errorf("user %d not found: %w", req.UserId, err)
	}
	// Получение задач по пользователю
	tasks, err := h.svc.ListTasksByUser(req.UserId, int(req.Page), int(req.PageSize))
	if err != nil {
		log.Printf("ListTasksByUser failed: %v", err)
		return nil, err
	}
	pbTasks := make([]*taskpb.Task, len(tasks))
	for i, t := range tasks {
		pbTasks[i] = &taskpb.Task{
			Id:     uint32(t.ID),
			UserId: t.UserID,
			Title:  t.Title,
			IsDone: t.IsDone,
		}
	}
	return &taskpb.ListTasksByUserResponse{Tasks: pbTasks}, nil
}
