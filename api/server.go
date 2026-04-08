package api

import (
	"database/sql"

	middleware "taskmanager/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const limit = 5
const offset = 0

func Start(db *sql.DB) {
	r := gin.Default()

	r.Use(middleware.RateLimiter())

	registerRoutes(r, db)

	r.Run(":8080")
}

func registerRoutes(r *gin.Engine, db *sql.DB) {

	//health
	r.GET("/health", HealthHandler(db))

	r.POST("/login", LoginHandler(db))

	r.POST("/register", RegisterHandler(db))

	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	authorized.GET("/tasks", GetTasksHandler(db))
	authorized.POST("/tasks", CreateTaskHandler(db))
	authorized.DELETE("/tasks/:id", DeleteTaskHandler(db))
	authorized.PATCH("/tasks/:id", UpdateTaskStatusHandler(db))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
