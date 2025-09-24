package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/I-Van-Radkov/kaspersky_1/internal/config"
	apphttp "github.com/I-Van-Radkov/kaspersky_1/internal/http"
	"github.com/I-Van-Radkov/kaspersky_1/internal/http/handlers"
	"github.com/I-Van-Radkov/kaspersky_1/internal/service"
)

type QueueServiceShutdownProvider interface {
	Shutdown()
	WaitForCompletion()
}

type App struct {
	HTTPServer   *apphttp.Server
	QueueService QueueServiceShutdownProvider
	done         chan struct{}
	shutdownOnce sync.Once
	wg           sync.WaitGroup
}

func New(cfg *config.Config) *App {
	queueService := service.NewQueueService(cfg.WorkerPool.Workers, cfg.WorkerPool.QueueSize)
	enqueueHandlers := handlers.NewEnqueueHandlers(queueService)
	httpRouter := apphttp.NewRouterGin(enqueueHandlers)

	return &App{
		HTTPServer:   apphttp.NewServer(cfg.HTTP.Port, cfg.HTTP.ReadTimeout, cfg.HTTP.WriteTimeout, httpRouter),
		QueueService: queueService,
		done:         make(chan struct{}),
	}
}

func (a *App) gracefulShutdown() {
	a.shutdownOnce.Do(func() {
		log.Println("Начинаем graceful shutdown")

		timeout := 30 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		log.Println("1. Останавливаем HTTP сервер")
		if err := a.HTTPServer.GracefulShutdown(5 * time.Second); err != nil {
			log.Printf("Ошибка остановки HTTP сервера: %v", err)
		}

		log.Println("2. Прекращаем прием новых задач")
		a.QueueService.Shutdown()

		log.Println("3. Ожидаем завершения текущих задач")

		completionDone := make(chan struct{})
		go func() {
			a.QueueService.WaitForCompletion()
			close(completionDone)
		}()

		select {
		case <-completionDone:
			log.Println("Все задачи завершены")
		case <-ctx.Done():
			log.Println("Таймаут ожидания завершения задач")
		}

		log.Println("Приложение корректно завершено")
		close(a.done)
	})
}

func (a *App) Run() {
	// Канал для ошибок сервера
	serverErr := make(chan error, 1)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		log.Println("Запуск HTTP сервера")
		if err := a.HTTPServer.Run(); err != nil {
			serverErr <- err
		}
	}()

	go a.setupSignalHandler()

	log.Println("Приложение запущено")

	select {
	case err := <-serverErr:
		log.Printf("Ошибка сервера: %v", err)
		a.gracefulShutdown()
	case <-a.done:
		// Graceful shutdown уже выполняется
	}

	a.wg.Wait()
	log.Println("Приложение полностью завершено")
}

func (a *App) setupSignalHandler() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Получен сигнал: %v", sig)
	a.gracefulShutdown()
}
