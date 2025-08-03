package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheck interface {
	HealthHandler(c *gin.Context)
}

func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": "1.0.0",
	})
}
