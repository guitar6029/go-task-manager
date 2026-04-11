package api

import (
	"database/sql"
	"log"

	middleware "taskmanager/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Start(db *sql.DB, rdb *redis.Client) {
	r := gin.Default()

	r.Use(middleware.RateLimiter())

	registerRoutes(r, db, rdb)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func registerRoutes(r *gin.Engine, db *sql.DB, rdb *redis.Client) {

	//health
	r.GET("/health", HealthHandler(db))

	r.POST("/login", LoginHandler(db))

	r.POST("/register", RegisterHandler(db))

	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	authorized.GET("/tasks", GetTasksHandler(db, rdb))
	authorized.POST("/tasks", CreateTaskHandler(db, rdb))
	authorized.DELETE("/tasks/:id", DeleteTaskHandler(db, rdb))
	authorized.PATCH("/tasks/:id", UpdateTaskStatusHandler(db, rdb))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
