package storage

import (
	"errors"

	"github.com/LNMMusic/tester-example/internal/task"
)

var (
	// ErrTaskNotFound is returned when the task is not found
	ErrTaskNotFound = errors.New("storage: task not found")
)

// Task is an interface that represents a storage
type Task interface {
	// FindById returns a task by its id
	FindById(id string) (t task.Task, err error)

	// Save saves a task
	Save(t *task.Task) (err error)
}