package main

import (
	"log"
	"taskmanager/internal/queue"

	redispkg "taskmanager/internal/redis"
)

func main() {

	ctx := redispkg.Ctx

	client := redispkg.NewClient()

	q := queue.NewRedisQueue(client, "jobs")

	log.Println("Worker started... waiting for jobs")

	for {
		job, err := q.PopJob(ctx)
		if err != nil {
			log.Println("error:", err)
		}
		log.Println("received job: ", job.Type)
	}
}
