package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "taskmanager/docs"
	"taskmanager/internal/api"

	dbpkg "taskmanager/internal/db"
	redispkg "taskmanager/internal/redis"

	envpkg "taskmanager/internal/config"
)

func main() {
	// Redis
	rdb := redispkg.NewClient()

	_, err := rdb.Ping(redispkg.Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis!")

	// load env
	envpkg.LoadEnv()

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
	api.Start(db, rdb)
}

func dbInit() (*sql.DB, error) {
	db, err := dbpkg.Init()
	if err != nil {
		return nil, fmt.Errorf("error initializing DB: %s", err)
	}
	return db, nil
}
