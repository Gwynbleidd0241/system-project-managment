package app

import (
	log "github.com/go-ozzo/ozzo-log"
	grpcApp "taskService/internal/app/grpc"
	"taskService/internal/service"
)

type App struct {
	GRPCServer *grpcApp.App
}

func New(log *log.Logger, grpcPort string, taskService *service.TaskService) *App {

	grpcApps := grpcApp.New(log, grpcPort, taskService)

	return &App{
		GRPCServer: grpcApps,
	}
}
