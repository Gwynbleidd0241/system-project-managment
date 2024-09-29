package auth

import (
	"context"
	"fmt"
	authv1 "github.com/Dragzet/gRPCProtosv2/gen/go/task"
	jwt "github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"taskService/internal/domain/models"
	"time"
)

const SecretKey = "secretKEY"

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
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	tasks, err := s.taskService.FindAllTasks(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "find all tasks error")
	}
	var protoTasks []*authv1.Task

	for _, task := range tasks {
		protoTask := &authv1.Task{
			Id:          task.ID,
			Email:       task.UserEmail,
			Name:        task.Name,
			Description: task.Description,
			Done:        task.IsDone,
		}
		protoTasks = append(protoTasks, protoTask)
	}

	return &authv1.GetTasksResponse{Tasks: protoTasks}, nil
}

func (s *ServerAPI) UpdateTask(ctx context.Context, req *authv1.UpdateTaskRequest) (*authv1.UpdateTaskResponse, error) {
	if req.Name == "" || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "name and email is required")
	}

	_, err := ValidToken(req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	task := models.Task{
		UserEmail:   req.Email,
		Name:        req.Name,
		Description: req.Description,
		IsDone:      req.Done,
	}

	err = s.taskService.UpdateTask(ctx, task)
	if err != nil {
		return &authv1.UpdateTaskResponse{Success: false}, status.Error(codes.Internal, "update task error")
	}
	return &authv1.UpdateTaskResponse{Success: true}, nil
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

	_, err := ValidToken(req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
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

func ValidToken(token string) (bool, error) {
	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		fmt.Println(err)
		return false, status.Error(codes.Unauthenticated, "invalid token")
	}

	claims := tokenParsed.Claims.(jwt.MapClaims)
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return false, status.Error(codes.Unauthenticated, "token expired")
		}
	}
	return true, nil

}
