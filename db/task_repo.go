package db

import (
	"database/sql"
	"fmt"
	model "taskmanager/model"
)

func CreateTask(db *sql.DB, title string) (int64, error) {
	if title == "" {
		return 0, fmt.Errorf("title cannot be empty")
	}
	result, err := db.Exec(`INSERT INTO tasks (title, done) VALUES (? , ?)`, title, false)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, err
}
func GetTasks(db *sql.DB, listType string, limit int, offset int) ([]model.Task, error) {
	//maybe set limit later
	var tasks = make([]model.Task, 0)
	query := "SELECT id, title, done FROM tasks"
	args := []interface{}{}

	switch listType {
	case "done":
		query += " WHERE done = ?"
		args = append(args, true)
	case "pending":
		query += " WHERE done = ?"
		args = append(args, false)
	}

	// pagination
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Done)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func UpdateTaskStatus(db *sql.DB, taskID int, done bool) error {
	result, err := db.Exec(`UPDATE tasks SET done = ? WHERE id = ?`, done, taskID)
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

func DeleteTask(db *sql.DB, taskID int) error {
	result, err := db.Exec(`DELETE FROM tasks WHERE id = ?`, taskID)
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
