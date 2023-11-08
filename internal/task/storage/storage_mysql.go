package storage

import (
	"database/sql"
	"errors"

	"github.com/LNMMusic/tester-example/internal/task"
)

// TaskMySQL is an implementation of Task interface for MySQL
type TaskMySQL struct {
	// db is the database connection
	db *sql.DB
}

// NewTaskMySQL returns a new TaskMySQL
func NewTaskMySQL(db *sql.DB) *TaskMySQL {
	return &TaskMySQL{db: db}
}

// FindById returns a task by its id
func (st *TaskMySQL) FindById(id string) (t task.Task, err error) {
	// query
	query := "SELECT `id`, `title`, `description`, `done` FROM `tasks` WHERE `id` = ?"

	// execute query
	err = st.db.QueryRow(query, id).Scan(&t.Id, &t.Title, &t.Description, &t.Done)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrTaskNotFound
			return
		}
		return
	}
	return
}

// Save saves a task
func (st *TaskMySQL) Save(t *task.Task) (err error) {
	// query
	query := "INSERT INTO `tasks` (`title`, `description`) VALUES (?, ?)"

	// execute query
	result, err := st.db.Exec(query, t.Title, t.Description)
	if err != nil {
		return
	}

	// check affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}
	if rowsAffected != 1 {
		err = errors.New("storage: expected 1 row affected")
		return
	}

	// get last insert id
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return
	}

	// set id
	t.Id = int(lastInsertId)

	return
}