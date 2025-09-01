package grpc

import (
	"context"

	taskpb "github.com/AleksKAG/project-protos/proto/task"
	"github.com/AleksKAG/tasks-service/internal/task"
)

type Handler struct {
	svc *task.Service
	taskpb.UnimplementedTaskServiceServer
}

func NewHandler(svc *task.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	t, err := h.svc.CreateTask(task.Task{Title: req.Title})
	if err != nil {
		return nil, err
	}
	return &taskpb.CreateTaskResponse{
		Task: &taskpb.Task{
			Id:    uint32(t.ID),
			Title: t.Title,
		},
	}, nil
}

func (h *Handler) GetTask(ctx context.Context, req *taskpb.Task) (*taskpb.Task, error) {
	t, err := h.svc.GetTask(req.Id)
	if err != nil {
		return nil, err
	}
	return &taskpb.Task{
		Id:    uint32(t.ID),
		Title: t.Title,
	}, nil
}

func (h *Handler) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.UpdateTaskResponse, error) {
	t, err := h.svc.UpdateTask(task.Task{
		Title: req.Title,
	})
	if err != nil {
		return nil, err
	}
	t.ID = uint(req.Id) // Устанавливаем ID для обновления
	return &taskpb.UpdateTaskResponse{
		Task: &taskpb.Task{
			Id:    uint32(t.ID),
			Title: t.Title,
		},
	}, nil
}

func (h *Handler) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	err := h.svc.DeleteTask(req.Id)
	if err != nil {
		return nil, err
	}
	return &taskpb.DeleteTaskResponse{Success: true}, nil
}

func (h *Handler) ListTasks(ctx context.Context, req *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	tasks, err := h.svc.ListTasks(int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}
	pbTasks := make([]*taskpb.Task, len(tasks))
	for i, t := range tasks {
		pbTasks[i] = &taskpb.Task{
			Id:    uint32(t.ID),
			Title: t.Title,
		}
	}
	return &taskpb.ListTasksResponse{Tasks: pbTasks}, nil
}
