package server

import (
	"net/http"
	"notes-api/internal/notes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewRouter(database *mongo.Database) *gin.Engine {

	r := gin.Default() // default instance is created with logger and recovery middleware

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ok":     true,
			"status": "healthy",
		})
	})

	notes.RegisterRoutes(r, database) // register notes routes

	return r
}
