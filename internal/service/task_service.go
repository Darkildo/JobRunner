package service

import (
	"JobRunner/internal/domain"
	"math/rand"
	"time"
)

type TaskService struct {
	storage Storage
	queue   chan string
}

func NewTaskService(storage Storage) *TaskService {
	return &TaskService{
		storage: storage,
		queue:   make(chan string, 100),
	}
}

func (s *TaskService) CreateTask() *domain.Task {
	return s.storage.Create()
}

func (s *TaskService) Enqueue(taskID string) {
	s.queue <- taskID
}

func (s *TaskService) StartWorkers(n int) {
	for i := 0; i < n; i++ {
		go s.worker()
	}
}

func (s *TaskService) worker() {
	for taskID := range s.queue {
		task, err := s.storage.Get(taskID)
		if err != nil {
			continue
		}

		// Обновляем статус задачи
		now := time.Now()
		task.State = domain.StateRunning
		task.StartedAt = &now
		s.storage.Update(task)

		// Симуляция I/O bound операции (3-5 минут)
		delay := time.Duration(180+rand.Intn(120)) * time.Second
		select {
		case <-time.After(delay):
			task.State = domain.StateCompleted
			task.Result = "success"
			now = time.Now()
			task.FinishedAt = &now
			s.storage.Update(task)
		case <-task.CancelChan:
			task.State = domain.StateCanceled
			now = time.Now()
			task.FinishedAt = &now
			s.storage.Update(task)
		}
	}
}

func (s *TaskService) GetTask(id string) (*domain.Task, error) {
	return s.storage.Get(id)
}

func (s *TaskService) CancelTask(id string) error {
	task, err := s.storage.Get(id)
	if err != nil {
		return err
	}

	if task.State == domain.StatePending || task.State == domain.StateRunning {
		close(task.CancelChan)
		task.State = domain.StateCanceled
		now := time.Now()
		task.FinishedAt = &now
		s.storage.Update(task)
	}

	return nil
}

func (s *TaskService) DeleteTask(id string) error {
	return s.storage.Delete(id)
}

func (s *TaskService) Shutdown() {
	close(s.queue)

	// Костыль чтобы воркеры завершились, по хорошему gracefully shutdown иначе реализуется
	time.Sleep(1 * time.Second)
}
