package main

import (
	"context"
	"encoding/json"
	"log"
	"taskmanager/internal/cache"
	"taskmanager/internal/queue"

	envpkg "taskmanager/internal/config"
	dbpkg "taskmanager/internal/db"
	redispkg "taskmanager/internal/redis"
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

	// DB
	db, err := dbpkg.Init()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()
	log.Println("worker connected to DB")

	// worker loop
	ctx := context.Background()
	log.Println("worker started... waiting for jobs")

	for {
		job, err := q.PopJob(ctx)
		if err != nil {
			log.Println("error:", err)
		}
		switch job.Type {

		case "create_task":
			var payload struct {
				Title string `json:"title"`
			}

			if err := json.Unmarshal(job.Payload, &payload); err != nil {
				log.Println("failed to parse payload:", err)
				continue
			}

			if _, err = dbpkg.CreateTask(db, payload.Title); err != nil {
				log.Println("failed to create task:", err)
				continue
			}

			cache.InvalidateTasks(rdb)
			log.Println("task created:", payload.Title)

		case "delete_task":
			var payload struct {
				ID int `json:"id"`
			}

			if err := json.Unmarshal(job.Payload, &payload); err != nil {
				log.Println("failed to parse payload:", err)
				continue
			}

			if err := dbpkg.DeleteTask(db, payload.ID); err != nil {
				log.Println("failed to delete task:", err)
				continue

			}

			cache.InvalidateTasks(rdb)
			log.Printf("Task %d deleted", payload.ID)

		case "mark_task_done":
			var payload struct {
				ID int `json:"id"`
			}

			if err := json.Unmarshal(job.Payload, &payload); err != nil {
				log.Println("failed to parse payload:", err)
				continue
			}

			if _, err := dbpkg.UpdateTaskStatus(db, payload.ID, true); err != nil {
				log.Println("failed to update task:", err)
				continue
			}

			cache.InvalidateTasks(rdb)
			log.Printf("Task id - %d has been updated to done", payload.ID)

		default:
			log.Println("unknown job type:", job.Type)
		}
	}
}
