package auth

import (
	"context"
	"fmt"
	apiv1 "github.com/Dragzet/gRPCAuthProtos/gen/go/auth"
	log "github.com/go-ozzo/ozzo-log"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Client struct {
	api    apiv1.AuthClient
	logger *log.Logger
}

func New(ctx context.Context, logger *log.Logger, addr string, timeout time.Duration, retriesCount int) (*Client, error) {
	const op = "grpc.New"

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

	return &Client{
		api:    apiv1.NewAuthClient(cc),
		logger: logger,
	}, nil
}

func interseptLogger(logger *log.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, args ...interface{}) {
		logger.Info(msg)
	})
}

func (c *Client) Register(ctx context.Context, username, password string) (string, error) {
	const op = "grpc.Auth"

	fmt.Println(username, password)

	resp, err := c.api.Register(ctx, &apiv1.RegisterRequest{
		Email:    username,
		Password: password,
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Token, nil
}

func (c *Client) Login(ctx context.Context, username, password string) (string, error) {
	const op = "grpc.Auth"

	resp, err := c.api.Login(ctx, &apiv1.LoginRequest{
		Email:    username,
		Password: password,
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Token, nil
}
