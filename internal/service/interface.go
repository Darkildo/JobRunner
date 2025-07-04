package service

import "JobRunner/internal/domain"

type Storage interface {
	Create() *domain.Task
	Get(id string) (*domain.Task, error)
	Update(task *domain.Task)
	Delete(id string) error
}
