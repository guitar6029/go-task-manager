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
func CreateTask(q *queue.RedisQueue, title string) error {
	if title == "" {
		return fmt.Errorf("title cannot be empty")
	}
	payload, err := json.Marshal(struct {
		Title string `json:"title"`
	}{Title: title})
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

func DeleteTask(id int, q *queue.RedisQueue) error {

	if id <= 0 {
		return fmt.Errorf("invalid id")
	}

	payload, err := json.Marshal(struct {
		ID int `json:"id"`
	}{ID: id})
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

func MarkTaskDone(id int, q *queue.RedisQueue) error {

	if id <= 0 {
		return fmt.Errorf("invalid id")
	}

	payload, err := json.Marshal(struct {
		ID int `json:"id"`
	}{ID: id})

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
