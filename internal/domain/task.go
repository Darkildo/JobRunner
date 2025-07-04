package domain

import (
	"fmt"
	"time"
)

type TaskState string

const (
	StatePending   TaskState = "pending"
	StateRunning   TaskState = "running"
	StateCompleted TaskState = "completed"
	StateFailed    TaskState = "failed"
	StateCanceled  TaskState = "canceled"
)

type Task struct {
	ID         string
	State      TaskState
	CreatedAt  time.Time
	StartedAt  *time.Time
	FinishedAt *time.Time
	Result     string
	CancelChan chan struct{}
}

func (t *Task) Duration() float64 {
	if t.StartedAt == nil {
		return 0
	}

	endTime := time.Now()
	if t.FinishedAt != nil {
		endTime = *t.FinishedAt
	}

	return endTime.Sub(*t.StartedAt).Seconds()
}

var ErrTaskNotFound = fmt.Errorf("task no found")
