package repository

import (
	"context"
	"fmt"
	"taskService/internal/domain/models"
	PostgreSQL "taskService/internal/storage/postgreSQL"
)

type StorageClient struct {
	Client PostgreSQL.Client
}

func (s *StorageClient) FindAllTasks(ctx context.Context, email string) ([]models.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (s *StorageClient) DeleteTask(ctx context.Context, taskName string) error {
	const taskDBErrorStatement = "internal/repository/tasks.go/DeleteTask: "
	stmt := `DELETE FROM tasks WHERE name = $1`

	_, err := s.FindTaskByName(ctx, taskName)
	if err != nil {
		return fmt.Errorf("no task with this name")
	}

	_, err = s.Client.Exec(ctx, stmt, taskName)
	if err != nil {
		return fmt.Errorf(taskDBErrorStatement + err.Error())
	}
	return nil
}

func (s *StorageClient) CreateTask(ctx context.Context, task models.Task) error {
	const taskDBErrorStatement = "internal/repository/tasks.go/NewTask: "

	stmt := `INSERT INTO tasks (user_email, name, description) VALUES ($1, $2, $3)`

	_, err := s.FindTaskByName(ctx, task.Name)
	if err == nil {
		fmt.Println(err)
		return fmt.Errorf("task already exists")
	}

	s.Client.QueryRow(ctx, stmt, task.UserEmail, task.Name, task.Description)
	return nil
}

func (s *StorageClient) FindTaskByName(ctx context.Context, name string) (models.Task, error) {
	const taskDBErrorStatement = "internal/repository/tasks.go/FindTaskByName: "

	stmt := `SELECT * FROM tasks WHERE name = $1`
	var task models.Task

	err := s.Client.QueryRow(ctx, stmt, name).Scan(&task.ID, &task.UserEmail, &task.Name, &task.Description)
	if err != nil {
		return task, fmt.Errorf("%s : %s", taskDBErrorStatement, err.Error())
	}
	return task, nil
}

func (s *StorageClient) UpdateTask(ctx context.Context, task models.Task) error {
	const taskDBErrorStatement = "internal/repository/tasks.go/UpdateTask: "

	stmt := `UPDATE tasks SET name = $1, description = $2 WHERE name = $3`

	_, err := s.FindTaskByName(ctx, task.Name)
	if err != nil {
		return fmt.Errorf("%s : %s", taskDBErrorStatement, err)
	}

	s.Client.QueryRow(ctx, stmt, task.Name, task.Description, task.ID)
	return nil
}

func New(client PostgreSQL.Client) *StorageClient {
	return &StorageClient{
		Client: client,
	}
}
