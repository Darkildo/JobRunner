package transport

import (
	"JobRunner/internal/domain"
)

type Service interface {
	CreateTask() *domain.Task
	Enqueue(taskID string)
	StartWorkers(n int)
	GetTask(id string) (*domain.Task, error)
	CancelTask(id string) error
	DeleteTask(id string) error
}
