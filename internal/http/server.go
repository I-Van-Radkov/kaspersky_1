package http

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	port         int
	readTimeout  time.Duration
	writeTimeout time.Duration
	router       http.Handler
	server       *http.Server
}

func NewServer(port int, readTimeout, writeTimeout time.Duration, router http.Handler) *Server {
	return &Server{
		port:         port,
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		router:       router,
	}
}

func (a *Server) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *Server) Run() error {
	a.server = &http.Server{
		Addr:         fmt.Sprintf(":%v", a.port),
		Handler:      a.router,
		ReadTimeout:  a.readTimeout,
		WriteTimeout: a.writeTimeout,
	}

	fmt.Printf("Сервер запущен на порту %d\n", a.port)
	err := a.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *Server) GracefulShutdown(timeout time.Duration) error {
	if a.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Println("Остановка HTTP сервера")
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("ошибка остановки HTTP сервера: %v", err)
	}

	fmt.Println("HTTP сервер корректно остановлен")
	return nil
}
