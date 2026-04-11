package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func HealthHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(500, gin.H{"status": "db down"})
		}
		c.JSON(200, gin.H{"status": "ok"})
	}
}
