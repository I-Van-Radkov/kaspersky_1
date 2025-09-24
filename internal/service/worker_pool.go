package service

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/I-Van-Radkov/kaspersky_1/internal/models"
)

type WorkerPool interface {
	Start()
	Shutdown()
	AddTask(task *models.Task)
	SetStatusTask(id string, status statusTask)
	WaitForCompletion()
	IsRunning() bool
}

type Pool struct {
	tasksStatus    map[string]statusTask
	tasks          chan *models.Task
	wg             sync.WaitGroup
	startOnce      sync.Once
	shutdownOnce   sync.Once
	mutex          sync.Mutex
	numWorkers     int
	isRunning      bool
	isShuttingDown bool
}

func NewPool(numWorkers, queueSize int) *Pool {
	p := &Pool{
		tasksStatus: make(map[string]statusTask),
		tasks:       make(chan *models.Task, queueSize),
		numWorkers:  numWorkers,
	}

	p.Start()

	return p
}

func (p *Pool) Start() {
	p.mutex.Lock()
	if p.isRunning {
		p.mutex.Unlock()
		return
	}

	p.mutex.Unlock()

	p.startOnce.Do(func() {
		p.mutex.Lock()
		p.isRunning = true
		p.isShuttingDown = false
		p.mutex.Unlock()

		for i := 0; i < p.numWorkers; i++ {
			p.wg.Add(1)
			go func() {
				p.work(i)
			}()
		}

		log.Println("WorkerPool запущен")
	})
}

func (p *Pool) work(id int) {
	defer p.wg.Done()

	for task := range p.tasks {
		err := p.processTask(task)
		if err != nil {
			log.Printf("Worker %d: %v\n", id, err)
		}
	}
}

func (p *Pool) processTask(task *models.Task) error {
	backoff := NewBackoffWithJitter(task.MaxRetries)

	for {
		log.Printf("Идет обработка задания id: %s", task.Id)

		p.SetStatusTask(task.Id, statusRunning) // обновление состояния

		processingTime := time.Duration(100+rand.Intn(400)) * time.Millisecond
		time.Sleep(processingTime) // симуляция обработки

		if rand.Float64() < 0.2 {
			delay, hasNext := backoff.Next()
			if !hasNext {
				log.Printf("Задача %s провалилась после %d попыток", task.Id, task.MaxRetries)

				p.SetStatusTask(task.Id, statusFailed) // обновление состояния
				return fmt.Errorf("задача %s провалилась после %d попыток", task.Id, task.MaxRetries)
			}

			log.Printf("Задача %s: повтор через %v\n", task.Id, delay)
			time.Sleep(delay)
			continue
		}

		log.Printf("Задача %s: успешно\n", task.Id)
		p.SetStatusTask(task.Id, statusDone) // обновление состояния
		return nil
	}
}

func (p *Pool) Shutdown() {
	p.shutdownOnce.Do(func() {
		log.Println("WorkerPool: начало graceful shutdown")

		p.mutex.Lock()
		p.isShuttingDown = true // Прекращаем принимать новые задачи
		p.mutex.Unlock()

		close(p.tasks)

		// Ждем завершения текущих задач
		log.Println("WorkerPool: ожидаем завершения текущих задач")
		p.wg.Wait()

		p.mutex.Lock()
		p.isRunning = false
		p.mutex.Unlock()

		log.Println("WorkerPool: graceful shutdown завершен")
	})
}

func (p *Pool) WaitForCompletion() {
	p.wg.Wait()
}

func (p *Pool) AddTask(task *models.Task) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isShuttingDown {
		log.Printf("Задача %s отклонена: WorkerPool завершает работу", task.Id)
		return
	}

	if p.isRunning {
		log.Printf("Задача %s добавлена в канал\n", task.Id)
		p.tasks <- task
	} else {
		log.Printf("Задача %s отклонена: WorkerPool не запущен\n", task.Id)
	}
}

func (p *Pool) SetStatusTask(id string, status statusTask) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.tasksStatus[id] = status
}

func (p *Pool) IsRunning() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.isRunning && !p.isShuttingDown
}
