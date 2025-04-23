package server

import (
	"fmt"
	"net/http"

	"github.com/Edu4rdoNeves/ingestor-magalu/cmd/api/router"
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/configs/env"
	"github.com/gin-gonic/gin"
)

func Run() {
	server := gin.Default()
	router := router.Router(server)

	router.GET("/health", healthRoute)

	port := env.ServerPort
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}

func healthRoute(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Health Response Okay"})
}
