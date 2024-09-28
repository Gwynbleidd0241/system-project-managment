package PostgreSQL

import (
	"authService/internal/domain/models"
	"authService/internal/storage/PostgreSQL/PostgreSQLClient"
	"context"
	"fmt"
)

type StorageClient struct {
	Client PostgreSQLClient.Client
}

func (s StorageClient) FindUser(ctx context.Context, email string) (models.User, error) {
	const userDBErrorStatement = "internal/storage/PostgreSQL/storage.go: "

	stmt := `SELECT id, email, password FROM users WHERE email = $1`
	var user models.User

	err := s.Client.QueryRow(ctx, stmt, email).Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		return user, fmt.Errorf("%s find userModule: %s", userDBErrorStatement, err.Error())
	}
	return user, nil
}

func (s StorageClient) SaveUser(ctx context.Context, user models.User) (string, error) {
	const userDBErrorStatement = "internal/storage/PostgreSQL/storage.go: "

	stmt := `
		INSERT INTO users
			(email, password)
		VALUES
			($1, $2)
		RETURNING id
		`
	_, err := s.FindUser(ctx, user.Email)
	if err == nil {
		return "", fmt.Errorf("user already exists")
	}

	var id string

	if err := s.Client.QueryRow(ctx, stmt, user.Email, user.PassHash).Scan(&id); err != nil {
		return "", fmt.Errorf("%s create userModule: %s", userDBErrorStatement, err)
	}
	return id, nil
}

func New(client PostgreSQLClient.Client) *StorageClient {
	return &StorageClient{
		Client: client,
	}
}
