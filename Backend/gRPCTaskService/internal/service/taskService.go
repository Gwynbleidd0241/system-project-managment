package service

import (
	"context"
	"fmt"
	log "github.com/go-ozzo/ozzo-log"
	"taskService/internal/domain/models"
)

type TaskService struct {
	logger   *log.Logger
	provider TasksProvider
}

type TasksProvider interface {
	CreateTask(ctx context.Context, task models.Task) (string, error)
	FindAllTasks(ctx context.Context, email string) ([]models.Task, error)
	UpdateTask(ctx context.Context, task models.Task) error
	DeleteTask(ctx context.Context, taskName string) error
}

func (t *TaskService) CreateTask(ctx context.Context, task models.Task) (string, error) {
	const op = "TaskService.CreateTask"

	id, err := t.provider.CreateTask(ctx, task)
	if err != nil {
		t.logger.Error(op, err)
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (t *TaskService) DeleteTask(ctx context.Context, taskID string) error {
	const op = "TaskService.DeleteTask"

	if taskID == "" {
		return fmt.Errorf("%w: task name is empty", op)
	}

	err := t.provider.DeleteTask(ctx, taskID)
	if err != nil {
		t.logger.Error(op, err)
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (t *TaskService) UpdateTask(ctx context.Context, task models.Task) error {
	const op = "TaskService.UpdateTask"
	if task.Name == "" {
		return fmt.Errorf("%w: task name is empty", op)
	}
	err := t.provider.UpdateTask(ctx, task)
	if err != nil {
		t.logger.Error(op, err)
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (t *TaskService) FindAllTasks(ctx context.Context, email string) ([]models.Task, error) {
	const op = "TaskService.FindAllTasks"
	if email == "" {
		return nil, fmt.Errorf("%w: email is empty", op)
	}
	tasks, err := t.provider.FindAllTasks(ctx, email)
	if err != nil {
		t.logger.Error(op, err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tasks, nil
}

func New(logger *log.Logger, provider TasksProvider) *TaskService {
	return &TaskService{
		logger:   logger,
		provider: provider,
	}
}
