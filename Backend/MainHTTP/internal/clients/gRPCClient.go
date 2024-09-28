package clients

import (
	"context"
	log "github.com/go-ozzo/ozzo-log"
	"mainHTTP/internal/clients/auth"
	"mainHTTP/internal/clients/notification"
	"mainHTTP/internal/clients/task"
)

type GRPCClient struct {
	*auth.AuthClient
	*task.TaskClient
	*notification.NotificationClient
}

func NewGRPCClient(logger *log.Logger, authAddr string, taskAddr string, notifAddr string, retries int) *GRPCClient {

	authClient, err := auth.New(context.Background(), logger, authAddr, retries)
	if err != nil {
		panic(err)
	}

	taskClient, err := task.New(context.Background(), logger, taskAddr, retries)
	if err != nil {
		panic(err)
	}

	notificationClient, err := notification.New(context.Background(), logger, notifAddr, retries)
	if err != nil {
		panic(err)
	}

	return &GRPCClient{
		AuthClient:         authClient,
		TaskClient:         taskClient,
		NotificationClient: notificationClient,
	}
}
