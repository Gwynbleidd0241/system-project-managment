package main

import (
	"authService/internal/app"
	"authService/internal/config"
	"authService/internal/services/auth"
	"authService/internal/storage/PostgreSQL"
	"authService/internal/storage/PostgreSQL/PostgreSQLClient"
	"context"
	log "github.com/go-ozzo/ozzo-log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	cfg, err := config.LoadConfig()

	var wg sync.WaitGroup
	wg.Add(1)
	quit := make(chan struct{})
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	if err != nil {
		panic(err)
	}

	logger := setupLogger(cfg.Logs_path)
	logger.Open()
	defer logger.Close()

	logger.Info("Starting AuthService")

	storage, err := PostgreSQLClient.NewStorage(context.Background(), cfg.Storage_path)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	storageClient := PostgreSQL.New(storage)

	authService := auth.New(logger, cfg.TokenTTL, storageClient, storageClient)

	application := app.New(logger, cfg.GRPCConfig.GRPC_port, authService, cfg.TokenTTL)

	go startGRPCServer(logger, application, &wg, quit)

	<-stop
	close(quit)
	wg.Wait()
}

func setupLogger(logsPath string) *log.Logger {
	logger := log.NewLogger()

	t1 := log.NewConsoleTarget()
	t2 := log.NewFileTarget()
	t2.FileName = logsPath
	logger.Targets = append(logger.Targets, t1, t2)

	return logger
}

func startGRPCServer(logger *log.Logger, application *app.App, wg *sync.WaitGroup, quit chan struct{}) {
	defer wg.Done()

	go func() {
		<-quit
		logger.Info("Завершение gRPC сервера...")
		application.GRPCServer.Stop()
	}()
	application.GRPCServer.MustRun()
}

func startHTTPServer(logger *log.Logger, server *http.Server, wg *sync.WaitGroup, quit chan struct{}) {
	defer wg.Done()

	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Ошибка при завершении HTTP сервера: %v", err)
		}
	}()

	logger.Info("Запуск HTTP сервера на " + server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Ошибка HTTP сервера: ", err)
	}
}
