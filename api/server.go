package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
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
}
