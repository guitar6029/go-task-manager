package queue

import (
	"context"
	"encoding/json"
	"taskmanager/internal/model"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client *redis.Client
	key    string
}

func NewRedisQueue(client *redis.Client, key string) *RedisQueue {
	return &RedisQueue{
		client: client,
		key:    key,
	}
}

func (q *RedisQueue) PushJob(ctx context.Context, job model.Job) error {
	// step 1 ) convert job -> JSON
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	// step 2 ) push into redis list
	return q.client.LPush(ctx, q.key, data).Err()
}
func (q *RedisQueue) PopJob(ctx context.Context) (*model.Job, error) {
	// step 1 ) block and wait for the job
	result, err := q.client.BRPop(ctx, 0, q.key).Result()
	if err != nil {
		return nil, err
	}

	// step 2 ) extract job data
	data := result[1]

	// step 3 ) convert json to job struct
	var job model.Job
	if err := json.Unmarshal([]byte(data), &job); err != nil {
		return nil, err
	}

	// step 4 ) return job
	return &job, nil
}
