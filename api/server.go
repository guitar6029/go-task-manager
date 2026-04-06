package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const limit = 5
const offset = 0

func Start(db *sql.DB) {
	r := gin.Default()

	registerRoutes(r, db)

	r.Run(":8080")
}

func registerRoutes(r *gin.Engine, db *sql.DB) {
	r.GET("/tasks", GetTasksHandler(db))

	r.POST("/tasks", CreateTaskHandler(db))

	r.DELETE("/tasks/:id", DeleteTaskHandler(db))

	r.PATCH("/tasks/:id", UpdateTaskStatusHandler(db))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
