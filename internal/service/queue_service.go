package service

import (
	"log"

	"github.com/I-Van-Radkov/kaspersky_1/internal/models"
)

type statusTask string

const (
	statusQueued  statusTask = "queued"
	statusRunning statusTask = "running"
	statusDone    statusTask = "done"
	statusFailed  statusTask = "failed"
)

type QueueService struct {
	pool WorkerPool
}

func NewQueueService(numWorkers, queueSize int) *QueueService {
	pool := NewPool(numWorkers, queueSize)

	qs := &QueueService{
		pool: pool,
	}

	return qs
}

func (s *QueueService) AddToQueue(id, payload string, maxRetries int) {
	task := models.ToTask(id, payload, maxRetries)

	s.pool.AddTask(task)

	s.pool.SetStatusTask(id, statusQueued)
}

func (s *QueueService) Shutdown() {
	log.Println("QueueService: начало graceful shutdown")
	s.pool.Shutdown()
	log.Println("QueueService: graceful shutdown завершен")
}

func (s *QueueService) WaitForCompletion() {
	log.Println("QueueService: ожидаем завершения задач")
	s.pool.WaitForCompletion()
	log.Println("QueueService: все задачи завершены")
}

func (s *QueueService) IsRunning() bool {
	return s.pool.IsRunning()
}
