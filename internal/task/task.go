package task

// Task is an struct that represents a task
type Task struct {
	// Id is the unique identifier of the task
	Id int `json:"id"`
	// Title is the title of the task
	Title string `json:"title"`
	// Description is the description of the task
	Description string `json:"description"`
	// Done is the status of the task
	Done bool `json:"done"`
}