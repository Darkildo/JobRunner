package main

import (
	"JobRunner/internal/service"
	storage2 "JobRunner/internal/storage"
	http2 "JobRunner/internal/transport/http"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// типо DI
	storage := storage2.NewMemoryStorage()
	taskService := service.NewTaskService(storage)
	taskHandler := http2.NewTaskHandler(taskService)

	// Запуск воркеров
	taskService.StartWorkers(5)

	// Настройка HTTP сервера
	router := http2.NewRouter(taskHandler)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Перехватываем сигнал завершения работы из ОС и пытаемся завершиться
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on port 8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	// Завершение всех задач
	taskService.Shutdown()
	log.Println("Server gracefully stopped")
}
