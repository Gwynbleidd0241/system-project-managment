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
	const op = "internal/repository/tasks.go/FindAllTasks: "
	stmt := `SELECT id, user_email, name, description, is_done FROM tasks WHERE user_email = $1`

	// Выполнение запроса
	rows, err := s.Client.Query(ctx, stmt, email)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err) // Обработка ошибки
	}
	defer rows.Close() // Закрытие rows после завершения работы

	// Слайс для хранения задач
	var tasks []models.Task

	// Итерация по строкам результата
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.UserEmail, &task.Name, &task.Description, &task.IsDone); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err) // Обработка ошибки
		}
		tasks = append(tasks, task) // Добавление задачи в слайс
	}

	// Проверка на ошибки после итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err) // Обработка ошибки
	}

	return tasks, nil // Возврат найденных задач
}

func (s *StorageClient) DeleteTask(ctx context.Context, id string) error {
	const taskDBErrorStatement = "internal/repository/tasks.go/DeleteTask: "
	stmt := `DELETE FROM tasks WHERE id = $1`

	_, err := s.FindTaskByID(ctx, id)
	if err != nil {
		return fmt.Errorf("no task with this name")
	}

	_, err = s.Client.Exec(ctx, stmt, id)
	if err != nil {
		return fmt.Errorf(taskDBErrorStatement + err.Error())
	}
	return nil
}

func (s *StorageClient) CreateTask(ctx context.Context, task models.Task) (string, error) {
	const taskDBErrorStatement = "internal/repository/tasks.go/NewTask: "

	stmt := `INSERT INTO tasks (user_email, name, description) VALUES ($1, $2, $3) Returning id`

	_, err := s.FindTaskByName(ctx, task.Name)
	if err == nil {
		return "", fmt.Errorf("task already exists")
	}

	err = s.Client.QueryRow(ctx, stmt, task.UserEmail, task.Name, task.Description).Scan(&task.ID)
	if err != nil {
		return "", fmt.Errorf("%s : %s", taskDBErrorStatement, err.Error())
	}
	return task.ID, nil
}

func (s *StorageClient) FindTaskByName(ctx context.Context, name string) (models.Task, error) {
	const taskDBErrorStatement = "internal/repository/tasks.go/FindTaskByName: "

	stmt := `SELECT * FROM tasks WHERE name = $1`
	var task models.Task

	err := s.Client.QueryRow(ctx, stmt, name).Scan(&task.ID, &task.UserEmail, &task.Name, &task.Description, &task.IsDone)
	if err != nil {
		return task, fmt.Errorf("%s : %s", taskDBErrorStatement, err.Error())
	}
	return task, nil
}

func (s *StorageClient) FindTaskByID(ctx context.Context, id string) (models.Task, error) {
	const taskDBErrorStatement = "internal/repository/tasks.go/FindTaskByName: "

	stmt := `SELECT * FROM tasks WHERE id = $1`
	var task models.Task

	err := s.Client.QueryRow(ctx, stmt, id).Scan(&task.ID, &task.UserEmail, &task.Name, &task.Description, &task.IsDone)
	fmt.Println(task)
	if err != nil {
		return task, fmt.Errorf("%s : %s", taskDBErrorStatement, err.Error())
	}
	return task, nil
}

func (s *StorageClient) UpdateTask(ctx context.Context, task models.Task) error {
	const taskDBErrorStatement = "internal/repository/tasks.go/UpdateTask: "

	stmt := `UPDATE tasks SET description = $2, is_done = $3  WHERE name = $1`

	_, err := s.FindTaskByName(ctx, task.Name)
	if err != nil {
		return fmt.Errorf("%s : %s", taskDBErrorStatement, err)
	}

	_, err = s.Client.Exec(ctx, stmt, task.Name, task.Description, task.IsDone)
	if err != nil {
		return fmt.Errorf("%s : %s", taskDBErrorStatement, err.Error())
	}
	return nil
}

func New(client PostgreSQL.Client) *StorageClient {
	return &StorageClient{
		Client: client,
	}
}
