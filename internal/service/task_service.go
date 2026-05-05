package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"taskmanager/internal/cache"
	dbpkg "taskmanager/internal/db"
	model "taskmanager/internal/model"
	"taskmanager/internal/queue"

	"github.com/redis/go-redis/v9"
)

func GetTasks(db *sql.DB, rdb *redis.Client, userID int, filter string, limit int, offset int) ([]model.Task, error) {
	if limit <= 0 {
		limit = 5
	}

	cacheKey := fmt.Sprintf("tasks:%d:%s:%d:%d", userID, filter, limit, offset)

	//try cache
	tasks, found := cache.GetTasks(rdb, cacheKey)
	if found {
		return tasks, nil
	}

	return dbpkg.GetTasks(db, userID, filter, limit, offset)
}

func TaskBelongsToUser(db *sql.DB, taskID int, userID int) (bool, error) {
	if taskID <= 0 {
		return false, fmt.Errorf("invalid id")
	}

	return dbpkg.TaskBelongsToUser(db, taskID, userID)
}

func CreateTask(q *queue.RedisQueue, userID int, title string) error {
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	payload, err := json.Marshal(struct {
		Title  string `json:"title"`
		UserID int    `json:"user_id"`
	}{Title: title, UserID: userID})
	if err != nil {
		return err
	}

	job := model.Job{
		Type:     "create_task",
		Payload:  payload,
		Retries:  0,
		MaxRetry: 3,
	}

	return q.PushJob(context.Background(), job)
}

func DeleteTask(id int, userID int, q *queue.RedisQueue) error {

	if id <= 0 {
		return fmt.Errorf("invalid id")
	}

	payload, err := json.Marshal(struct {
		ID     int `json:"id"`
		UserID int `json:"user_id"`
	}{ID: id, UserID: userID})
	if err != nil {
		return err
	}

	job := model.Job{
		Type:     "delete_task",
		Payload:  payload,
		Retries:  0,
		MaxRetry: 3,
	}

	return q.PushJob(context.Background(), job)
}

func MarkTaskDone(id int, userID int, q *queue.RedisQueue) error {

	if id <= 0 {
		return fmt.Errorf("invalid id")
	}

	payload, err := json.Marshal(struct {
		ID     int `json:"id"`
		UserID int `json:"user_id"`
	}{ID: id, UserID: userID})

	if err != nil {
		return err
	}

	job := model.Job{
		Type:     "mark_task_done",
		Payload:  payload,
		Retries:  0,
		MaxRetry: 3,
	}

	return q.PushJob(context.Background(), job)
}
