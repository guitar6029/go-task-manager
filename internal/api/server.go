package api

import (
	"database/sql"
	"log"

	middleware "taskmanager/internal/middleware"
	"taskmanager/internal/queue"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Start(db *sql.DB, rdb *redis.Client, q *queue.RedisQueue) {
	r := gin.Default()

	r.Use(middleware.RateLimiter())

	registerRoutes(r, db, rdb, q)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func registerRoutes(r *gin.Engine, db *sql.DB, rdb *redis.Client, q *queue.RedisQueue) {

	//health
	r.GET("/health", HealthHandler(db))

	r.POST("/login", LoginHandler(db))

	r.POST("/register", RegisterHandler(db))

	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	authorized.GET("/tasks", GetTasksHandler(db, rdb))
	authorized.POST("/tasks", CreateTaskHandler(q))
	authorized.DELETE("/tasks/:id", DeleteTaskHandler(q))
	authorized.PATCH("/tasks/:id", UpdateTaskStatusHandler(q))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
