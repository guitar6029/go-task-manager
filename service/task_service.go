package service

import (
	"database/sql"
	"fmt"
	dbpkg "taskmanager/db"
	model "taskmanager/model"
)

func GetTasks(db *sql.DB, filter string, limit int, offset int) ([]model.Task, error) {
	if limit <= 0 {
		limit = 5
	}

	return dbpkg.GetTasks(db, filter, limit, offset)
}

func CreateTask(db *sql.DB, title string) (int64, error) {
	if title == "" {
		return 0, fmt.Errorf("title cannot be empty")
	}
	return dbpkg.CreateTask(db, title)
}

func DeleteTask(db *sql.DB, id int) error {
	return dbpkg.DeleteTask(db, id)
}

func MarkTaskDone(db *sql.DB, id int) (model.Task, error) {
	return dbpkg.UpdateTaskStatus(db, id, true)
}
