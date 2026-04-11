package service

import (
	"database/sql"
	"fmt"
	"taskmanager/internal/cache"
	dbpkg "taskmanager/internal/db"
	model "taskmanager/internal/model"

	"github.com/redis/go-redis/v9"
)

func GetTasks(db *sql.DB, rdb *redis.Client, filter string, limit int, offset int) ([]model.Task, error) {
	if limit <= 0 {
		limit = 5
	}

	cacheKey := fmt.Sprintf("tasks:%s:%d:%d", filter, limit, offset)

	//try cache
	tasks, found := cache.GetTasks(rdb, cacheKey)
	if found {
		return tasks, nil
	}

	return dbpkg.GetTasks(db, filter, limit, offset)
}

func CreateTask(db *sql.DB, rdb *redis.Client, title string) (int64, error) {
	if title == "" {
		return 0, fmt.Errorf("title cannot be empty")
	}
	id, err := dbpkg.CreateTask(db, title)
	if err != nil {
		return id, err
	}

	cache.InvalidateTasks(rdb)
	return id, nil
}

func DeleteTask(db *sql.DB, rdb *redis.Client, id int) error {
	err := dbpkg.DeleteTask(db, id)
	if err != nil {
		return err
	}

	cache.InvalidateTasks(rdb)

	return nil
}

func MarkTaskDone(db *sql.DB, rdb *redis.Client, id int) (model.Task, error) {
	task, err := dbpkg.UpdateTaskStatus(db, id, true)
	if err != nil {
		return model.Task{}, err
	}

	cache.InvalidateTasks(rdb)
	return task, nil
}
