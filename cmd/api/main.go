package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "taskmanager/docs/api"
	"taskmanager/internal/api"
	"taskmanager/internal/queue"

	dbpkg "taskmanager/internal/db"
	redispkg "taskmanager/internal/redis"

	envpkg "taskmanager/internal/config"
)

func main() {
	// load env
	envpkg.LoadEnv()

	// Redis
	rdb := redispkg.NewClient()

	_, err := rdb.Ping(redispkg.Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis!")

	// queue
	q := queue.NewRedisQueue(rdb, "jobs")

	// DB
	db, err := dbInit()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("error closing db: ", err)
		}
	}()

	// Start API
	fmt.Println("Initializing API Program")
	api.Start(db, rdb, q)
}

func dbInit() (*sql.DB, error) {
	db, err := dbpkg.Init()
	if err != nil {
		return nil, fmt.Errorf("error initializing DB: %s", err)
	}
	return db, nil
}
