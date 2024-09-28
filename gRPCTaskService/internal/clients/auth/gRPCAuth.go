package auth

//
//import (
//	"context"
//	"fmt"
//	apiv1 "github.com/Dragzet/gRPCAuthProtos/gen/go/auth"
//	log "github.com/go-ozzo/ozzo-log"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/credentials/insecure"
//	"time"
//)
//
//type Client struct {
//	api    apiv1.AuthClient
//	logger *log.Logger
//}
//
//func New(ctx context.Context, logger *log.Logger, addr string, timeout time.Duration, retriesCount int) (*Client, error) {
//	const op = "grpc.New"
//
//	cc, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
//	if err != nil {
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//
//}
