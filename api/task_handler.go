package api

import (
	"database/sql"
	"net/http"
	servicepkg "taskmanager/service"

	"github.com/gin-gonic/gin"
)

func GetTasksHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tasks, err := servicepkg.GetTasks(db, "", 5, 0)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tasks)
	}
}
