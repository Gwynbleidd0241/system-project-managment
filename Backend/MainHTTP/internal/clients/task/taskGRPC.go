package task

import (
	"context"
	"fmt"
	apiv1 "github.com/Dragzet/gRPCProtosv2/gen/go/task"
	log "github.com/go-ozzo/ozzo-log"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type TaskClient struct {
	api    apiv1.TaskServiceClient
	logger *log.Logger
}

func New(ctx context.Context, logger *log.Logger, addr string, retriesCount int) (*TaskClient, error) {
	const op = "grpc.New"
	addr = "task-service:" + addr

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithBackoff(grpcretry.BackoffLinear(time.Second)),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpcretry.UnaryClientInterceptor(retryOpts...),
			grpclog.UnaryClientInterceptor(interseptLogger(logger), logOpts...),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &TaskClient{
		api:    apiv1.NewTaskServiceClient(cc),
		logger: logger,
	}, nil
}

func interseptLogger(logger *log.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, args ...interface{}) {
		logger.Info(msg)
	})
}

func (c *TaskClient) CreateTask(ctx context.Context, task *apiv1.CreateTaskRequest) (string, error) {
	const op = "grpc.Task"

	resp, err := c.api.CreateTask(ctx, task)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Id, nil
}

func (c *TaskClient) DeleteTask(ctx context.Context, task *apiv1.DeleteTaskRequest) (bool, error) {
	const op = "grpc.Task"

	resp, err := c.api.DeleteTask(ctx, task)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return resp.Success, nil
}

func (c *TaskClient) UpdateTask(ctx context.Context, task *apiv1.UpdateTaskRequest) (bool, error) {
	const op = "grpc.Task"

	resp, err := c.api.UpdateTask(ctx, task)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return resp.Success, nil
}

func (c *TaskClient) GetTasks(ctx context.Context, email string) ([]*apiv1.Task, error) {
	const op = "grpc.Task"

	resp, err := c.api.GetTasks(ctx, &apiv1.GetTasksRequest{Email: email})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp.Tasks, nil
}
