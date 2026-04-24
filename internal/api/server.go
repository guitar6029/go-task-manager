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

// docker network range
var trustedProxies = []string{"172.18.0.0/16"}

func Start(db *sql.DB, rdb *redis.Client, q *queue.RedisQueue) {
	r := gin.Default()

	if err := r.SetTrustedProxies(trustedProxies); err != nil {
		log.Fatal("failed to set trusted proxies:", err)
	}

	r.Use(middleware.RateLimiter())

	registerRoutes(r, db, rdb, q)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func registerRoutes(r *gin.Engine, db *sql.DB, rdb *redis.Client, q *queue.RedisQueue) {

	//health
	r.GET("/health", HealthHandler(db))

	// debug ip
	r.GET("/debug-ip", func(c *gin.Context) {
		log.Println("Client IP:", c.ClientIP())
		c.JSON(200, gin.H{"ip": c.ClientIP()})
	})

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
