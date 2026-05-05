package db

import (
	"database/sql"

	"fmt"
	"log"
	model "taskmanager/internal/model"
)

// CreateTask inserts a new task and returns its ID
func CreateTask(db *sql.DB, title string) (int64, error) {
	if title == "" {
		return 0, fmt.Errorf("title cannot be empty")
	}

	var id int64
	err := db.QueryRow(
		`INSERT INTO tasks (title, done) VALUES ($1, $2) RETURNING id`,
		title,
		false,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetTasks retrieves tasks with optional filtering and pagination
func GetTasks(db *sql.DB, userID int, listType string, limit int, offset int) ([]model.Task, error) {

	var tasks = []model.Task{}

	query := "SELECT id, title, done FROM tasks WHERE user_id = $1"
	args := []interface{}{userID}
	argID := 2

	switch listType {
	case "done":
		query += fmt.Sprintf(" AND done = $%d", argID)
		args = append(args, true)
		argID++
	case "pending":
		query += fmt.Sprintf(" AND done = $%d", argID)
		args = append(args, false)
		argID++
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println("error closing rows: ", err)
		}
	}()

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Done); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// UpdateTaskStatus updates the done status of a task and returns the updated task
func UpdateTaskStatus(db *sql.DB, taskID int, done bool) (model.Task, error) {
	result, err := db.Exec(
		`UPDATE tasks SET done = $1 WHERE id = $2`,
		done,
		taskID,
	)
	if err != nil {
		return model.Task{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return model.Task{}, err
	}
	if rowsAffected == 0 {
		return model.Task{}, fmt.Errorf("task not found")
	}

	var task model.Task
	err = db.QueryRow(
		`SELECT id, title, done FROM tasks WHERE id = $1`,
		taskID,
	).Scan(&task.ID, &task.Title, &task.Done)

	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

// DeleteTask removes a task by ID
func DeleteTask(db *sql.DB, taskID int) error {
	result, err := db.Exec(
		`DELETE FROM tasks WHERE id = $1`,
		taskID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
