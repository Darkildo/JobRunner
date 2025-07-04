package storage

import (
	"JobRunner/internal/domain"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MemoryStorage struct {
	tasks sync.Map
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (s *MemoryStorage) Create() *domain.Task {
	task := &domain.Task{
		ID:         uuid.New().String(),
		State:      domain.StatePending,
		CreatedAt:  time.Now(),
		CancelChan: make(chan struct{}, 1),
	}
	s.tasks.Store(task.ID, task)
	return task
}

func (s *MemoryStorage) Get(id string) (*domain.Task, error) {
	task, ok := s.tasks.Load(id)
	if !ok {
		return nil, ErrTaskNotFound
	}
	return task.(*domain.Task), nil
}

func (s *MemoryStorage) Update(task *domain.Task) {
	s.tasks.Store(task.ID, task)
}

func (s *MemoryStorage) Delete(id string) error {
	s.tasks.Delete(id)
	return nil
}

var ErrTaskNotFound = fmt.Errorf("task not found")
