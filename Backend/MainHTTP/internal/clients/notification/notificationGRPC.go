package notification

import (
	"context"
	"fmt"
	apiv1 "github.com/Dragzet/gRPCProtosv2/gen/go/notification"
	log "github.com/go-ozzo/ozzo-log"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type NotificationClient struct {
	api    apiv1.MailerServiceClient
	logger *log.Logger
}

func New(ctx context.Context, logger *log.Logger, addr string, retriesCount int) (*NotificationClient, error) {
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

	return &NotificationClient{
		api:    apiv1.NewMailerServiceClient(cc),
		logger: logger,
	}, nil
}

func interseptLogger(logger *log.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, args ...interface{}) {
		logger.Info(msg)
	})
}

func (c *NotificationClient) SendEmail(ctx context.Context, email string) (string, error) {
	const op = "grpc.Notification"

	req := &apiv1.SendEmailRequest{
		To:          email,
		Subject:     "Welcome to our service",
		Body:        "Welcome to our service",
		Attachments: make([]string, 0),
	}

	resp, err := c.api.SendEmail(ctx, req)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.Message, nil
}
