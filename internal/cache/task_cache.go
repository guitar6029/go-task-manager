package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	model "taskmanager/internal/model"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func GetTasks(rdb *redis.Client, key string) ([]model.Task, bool) {
	cached, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, false
	}

	var tasks []model.Task
	if err := json.Unmarshal([]byte(cached), &tasks); err != nil {
		return nil, false
	}

	log.Println("Cache hit.")
	return tasks, true
}

func SetTasks(rdb *redis.Client, key string, tasks []model.Task) {
	data, err := json.Marshal(tasks)
	if err != nil {
		return
	}

	err = rdb.Set(ctx, key, data, time.Minute*5).Err()
	if err == nil {
		log.Println("Cache set.")
	}
}

func InvalidateTasks(rdb *redis.Client) {
	keys, err := rdb.Keys(ctx, "tasks:*").Result()
	if err != nil {
		log.Println("cache invalidation error:", err)
		return
	}

	if len(keys) > 0 {
		rdb.Del(ctx, keys...)
	}
}
