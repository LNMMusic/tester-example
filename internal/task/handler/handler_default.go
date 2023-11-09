package handler

import (
	"errors"
	"net/http"

	"github.com/LNMMusic/tester-example/internal/task"
	"github.com/LNMMusic/tester-example/internal/task/storage"
	"github.com/LNMMusic/tester-example/pkg/web/request"
	"github.com/LNMMusic/tester-example/pkg/web/response"
	"github.com/go-chi/chi/v5"
)

// Task represents a task handler
type Task struct {
	// st is the task storage
	st storage.Task
}

// NewTask returns a new Task
func NewTask(st storage.Task) *Task {
	return &Task{st: st}
}

// TaskJSON is an struct that represents a task in JSON format
type TaskJSON struct {
	// Id is the unique identifier of the task
	Id int `json:"id"`
	// Title is the title of the task
	Title string `json:"title"`
	// Description is the description of the task
	Description string `json:"description"`
	// Done is the status of the task
	Done bool `json:"done"`
}

// GetById returns a task by its id
func (hd *Task) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")

		// process
		t, err := hd.st.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, storage.ErrTaskNotFound):
				response.Error(w, http.StatusNotFound, "task not found")
			default:
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		// - serialize
		data := TaskJSON{
			Id:          t.Id,
			Title:       t.Title,
			Description: t.Description,
			Done:        t.Done,
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "task found",
			"data":    data,
		})
	}
}

// RequestBodyCreateTask is an struct that represents the request body for create task
type RequestBodyCreateTask struct {
	// Title is the title of the task
	Title string `json:"title"`
	// Description is the description of the task
	Description string `json:"description"`
}

// Create creates a task
func (hd *Task) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		var body RequestBodyCreateTask
		err := request.JSON(r, &body)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// process
		// - deserialize
		t := task.Task{
			Title:       body.Title,
			Description: body.Description,
		}
		err = hd.st.Save(&t)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// response
		// - serialize
		data := TaskJSON{
			Id:          t.Id,
			Title:       t.Title,
			Description: t.Description,
			Done:        t.Done,
		}
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "task created",
			"data":    data,
		})
	}
}