package app

import (
	grpcApp "authService/internal/app/grpc"
	"authService/internal/services/auth"
	log "github.com/go-ozzo/ozzo-log"
	"time"
)

type App struct {
	GRPCServer *grpcApp.App
}

func New(log *log.Logger, grpcPort string, authService *auth.Auth, ttl time.Duration) *App {

	grpcApps := grpcApp.New(log, grpcPort, authService)

	return &App{
		GRPCServer: grpcApps,
	}
}
