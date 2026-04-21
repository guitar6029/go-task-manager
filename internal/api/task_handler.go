package api

import (
	"database/sql"
	"net/http"
	"strconv"
	model "taskmanager/internal/model"
	"taskmanager/internal/queue"
	servicepkg "taskmanager/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var _ = model.Task{}

// GetTasks godoc
// @Summary Get tasks
// @Description Get tasks with optional filters
// @Tags tasks
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Param done query bool false "Done filter"
// @Success 200 {array} model.Task
// @Failure 500 {object} map[string]string
// @Router /tasks [get]
func GetTasksHandler(db *sql.DB, rdb *redis.Client) gin.HandlerFunc {
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

		tasks, err := servicepkg.GetTasks(db, rdb, filter, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tasks)
	}
}

type CreateTaskRequest struct {
	Title string `json:"title"`
}

// CreateTask godoc
// @Summary Create Task
// @Description Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body CreateTaskRequest true "Task payload"
// @Success 201 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func CreateTaskHandler(q *queue.RedisQueue) gin.HandlerFunc {
	return func(c *gin.Context) {

		var body CreateTaskRequest

		// parse JSON
		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "invalid body"})
			return
		}

		//validate it
		if body.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title cannot be empty"})
			return
		}

		//call service (queue now)
		err := servicepkg.CreateTask(q, body.Title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "task queued",
		})
	}
}

// DeleteTask godoc
// @Summary Delete Task
// @Description Delete a task
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 204 {string} string "No response body"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [delete]
func DeleteTaskHandler(q *queue.RedisQueue) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		err = servicepkg.DeleteTask(id, q)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}

		//c.Status(http.StatusNoContent) // 204
		c.JSON(http.StatusCreated, gin.H{
			"message": "task queued",
		})
	}
}

// UpdateTask godoc
// @Summary Update Task status
// @Description Update task's done property
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [patch]
func UpdateTaskStatusHandler(db *sql.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		task, err := servicepkg.MarkTaskDone(db, rdb, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}

		c.JSON(http.StatusOK, task)
	}
}
