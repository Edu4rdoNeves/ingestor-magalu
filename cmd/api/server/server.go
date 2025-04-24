package server

import (
	"fmt"
	"net/http"

	"github.com/Edu4rdoNeves/ingestor-magalu/cmd/api/router"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/configs/env"
	"github.com/gin-gonic/gin"
)

func Run() {
	port := env.ServerPort
	if port == "" {
		port = "8081"
	}

	server := gin.Default()
	server.GET("/health", healthRoute)

	router.SetupRouter(server)

	fmt.Printf("🚀 Server running on port %s\n", port)
	if err := server.Run(":" + port); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}

func healthRoute(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Health check passed!",
	})
}
