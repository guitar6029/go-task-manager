package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"taskmanager/internal/cache"
	"taskmanager/internal/model"
	"taskmanager/internal/queue"

	envpkg "taskmanager/internal/config"
	dbpkg "taskmanager/internal/db"
	redispkg "taskmanager/internal/redis"

	"github.com/redis/go-redis/v9"
)

func main() {

	// load env
	envpkg.LoadEnv()

	// redis
	rdb := redispkg.NewClient()
	_, err := rdb.Ping(redispkg.Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("worker connected to Redis")

	// queue
	q := queue.NewRedisQueue(rdb, "jobs")

	// failed queue
	qFailed := queue.NewRedisQueue(rdb, "jobs:failed")

	// DB
	db, err := dbpkg.Init()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("error closing db:", err)
		}
	}()

	log.Println("worker connected to DB")

	// worker loop
	ctx := context.Background()
	log.Println("worker started... waiting for jobs")

	for {
		job, err := q.PopJob(ctx)
		if err != nil {
			log.Println("pop error:", err)
			continue
		}

		log.Printf("processing job: %s (%d/%d)", job.Type, job.Retries, job.MaxRetry)

		if err := handleJob(db, rdb, job); err != nil {
			log.Println("job failed:", err)

			if job.Retries < job.MaxRetry {
				job.Retries++
				log.Printf("retrying job (%d/%d)", job.Retries, job.MaxRetry)

				if err := q.PushJob(ctx, *job); err != nil {
					log.Println("failed to requeue job:", err)
				}
			} else {
				log.Println("job permanently failed → moving to failed queue")

				if err := qFailed.PushJob(ctx, *job); err != nil {
					log.Println("failed to push to failed queue:", err)
				}
			}

			continue
		}

		log.Println("job succeeded:", job.Type)
	}
}

func handleJob(db *sql.DB, rdb *redis.Client, job *model.Job) error {
	switch job.Type {

	case "create_task":
		var payload struct {
			Title string `json:"title"`
		}

		if err := json.Unmarshal(job.Payload, &payload); err != nil {
			return err
		}

		if _, err := dbpkg.CreateTask(db, payload.Title); err != nil {

			return err
		}

		cache.InvalidateTasks(rdb)
		log.Println("task created:", payload.Title)
		return nil
		// for testing
		//return fmt.Errorf("forced failure for testing")

	case "delete_task":
		var payload struct {
			ID int `json:"id"`
		}

		if err := json.Unmarshal(job.Payload, &payload); err != nil {

			return err
		}

		if err := dbpkg.DeleteTask(db, payload.ID); err != nil {
			return err

		}

		cache.InvalidateTasks(rdb)
		log.Printf("Task %d deleted", payload.ID)
		return nil

	case "mark_task_done":
		var payload struct {
			ID int `json:"id"`
		}

		if err := json.Unmarshal(job.Payload, &payload); err != nil {
			return err
		}

		if _, err := dbpkg.UpdateTaskStatus(db, payload.ID, true); err != nil {
			return err
		}

		cache.InvalidateTasks(rdb)
		log.Printf("Task id - %d has been updated to done", payload.ID)
		return nil

	default:
		return fmt.Errorf("unknown job type: %s", job.Type)
	}
}
