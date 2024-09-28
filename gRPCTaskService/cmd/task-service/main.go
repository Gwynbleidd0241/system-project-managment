package main

import (
	"context"
	log "github.com/go-ozzo/ozzo-log"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"taskService/internal/config"
	"taskService/internal/http-server/handlers"
	MWLogger "taskService/internal/http-server/middleware/logger"
	"taskService/internal/repository"
	"taskService/internal/service"
	PostgreSQL "taskService/internal/storage/postgreSQL"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()

	var wg sync.WaitGroup
	wg.Add(2)
	quit := make(chan struct{})
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	if err != nil {
		panic(err) // No point in continuing if we can't load the cfg
	}

	logger := setupLogger(cfg.Logs_path)
	logger.Open()
	defer logger.Close()

	logger.Info("Starting task server")

	storage, err := PostgreSQL.NewStorage(context.Background(), cfg.Storage_path)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	storageClient := repository.New(storage)
	taskService := service.New(logger, storageClient)

	router := mux.NewRouter()
	router.Use(MWLogger.New(logger))
	handler := handlers.NewHandler(router, taskService)

	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go startHTTPServer(logger, server, &wg, quit)
	<-stop
	close(quit)
	wg.Wait()

}

func setupLogger(loggsPath string) *log.Logger {
	logger := log.NewLogger()

	t1 := log.NewConsoleTarget()
	t2 := log.NewFileTarget()
	t2.FileName = loggsPath
	logger.Targets = append(logger.Targets, t1, t2)

	return logger
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
