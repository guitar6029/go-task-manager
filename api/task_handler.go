package api

import (
	"database/sql"
	"net/http"
	"strconv"
	servicepkg "taskmanager/service"

	"github.com/gin-gonic/gin"
)

func GetTasksHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		//read query params
		limitStr := c.Query("limit")
		offsetStr := c.Query("offset")
		doneStr := c.Query("done")

		// convert
		limit := 5
		offset := 0
		filter := ""

		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}

		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}

		switch doneStr {
		case "true":
			filter = "done"
		case "false":
			filter = "pending"
		}

		tasks, err := servicepkg.GetTasks(db, filter, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tasks)
	}
}
