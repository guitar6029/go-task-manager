package api

import (
	"database/sql"
	"net/http"
	servicepkg "taskmanager/service"

	"github.com/gin-gonic/gin"
)

const limit = 5
const offset = 0

func Start(db *sql.DB) {
	r := gin.Default()

	r.GET("/tasks", func(c *gin.Context) {
		tasks, err := servicepkg.GetTasks(db, "", limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, tasks)
	})

	r.Run(":8080")
}
