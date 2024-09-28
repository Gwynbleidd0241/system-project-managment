package auth

import (
	"context"
	authv1 "github.com/Dragzet/gRPCProtos/gen/go/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"taskService/internal/domain/models"
)

type TaskServer interface {
	CreateTask(ctx context.Context, task models.Task) (string, error)
	FindAllTasks(ctx context.Context, email string) ([]models.Task, error)
	UpdateTask(ctx context.Context, task models.Task) error
	DeleteTask(ctx context.Context, taskName string) error
}

type ServerAPI struct {
	authv1.UnsafeTaskServiceServer
	taskService TaskServer
}

func RegisterServer(s *grpc.Server, taskServer TaskServer) {
	authv1.RegisterTaskServiceServer(s, &ServerAPI{
		taskService: taskServer,
	})
}

func (s *ServerAPI) GetTasks(ctx context.Context, req *authv1.GetTasksRequest) (*authv1.GetTasksResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServerAPI) UpdateTask(ctx context.Context, req *authv1.UpdateTaskRequest) (*authv1.UpdateTaskResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServerAPI) DeleteTask(ctx context.Context, req *authv1.DeleteTaskRequest) (*authv1.DeleteTaskResponse, error) {
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.taskService.DeleteTask(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "delete task error")
	}
	return &authv1.DeleteTaskResponse{Success: true}, nil
}

func (s *ServerAPI) CreateTask(ctx context.Context, req *authv1.CreateTaskRequest) (*authv1.CreateTaskResponse, error) {

	if req.Name == "" || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "name and email are required")
	}

	task := models.Task{
		UserEmail:   req.Email,
		Name:        req.Name,
		Description: req.Description,
	}

	id, err := s.taskService.CreateTask(ctx, task)
	if err != nil {
		return nil, status.Error(codes.Internal, "create task error")
	}
	return &authv1.CreateTaskResponse{Id: id}, nil
}
