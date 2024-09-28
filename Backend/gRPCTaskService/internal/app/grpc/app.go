package grpc

import (
	"fmt"
	log "github.com/go-ozzo/ozzo-log"
	"google.golang.org/grpc"
	"net"
	taskServer "taskService/internal/grpc/task"
)

type App struct {
	log        *log.Logger
	gRPCServer *grpc.Server
	port       string
}

func New(log *log.Logger, port string, taskService taskServer.TaskServer) *App {

	gRPCServer := grpc.NewServer()
	taskServer.RegisterServer(gRPCServer, taskService)
	return &App{
		gRPCServer: gRPCServer,
		log:        log,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		a.log.Error("Failed to start AuthService", "error", err)
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcApp.Run"
	a.log.Info("Starting gRPC server at port " + a.port)
	l, err := net.Listen("tcp", ":"+a.port)
	if err != nil {
		a.log.Error("failed to listen", "error", err, "operation", op)
		return fmt.Errorf("failed to listen", "error", err, "operation", op)
	}

	if err = a.gRPCServer.Serve(l); err != nil {
		a.log.Error("failed to serve", "error", err, "operation", op)
		return fmt.Errorf("failed to serve", "error", err, "operation", op)
	}
	return nil
}

func (a *App) Stop() {
	a.log.Info("Stopping gRPC server")
	a.gRPCServer.GracefulStop()
}
