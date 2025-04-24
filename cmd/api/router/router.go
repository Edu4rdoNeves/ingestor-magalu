package router

import (
	"net/http"
	"time"

	"github.com/Edu4rdoNeves/ingestor-magalu/internal/dependency"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {
	configureCORS(router)
	defineRoutes(router)
}

func configureCORS(router *gin.Engine) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNoContent)
	})
}

func defineRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		pulses := api.Group("/pulses")
		{
			pulses.POST("/populate", dependency.PulseController.PopulateQueueWithPulses)
			pulses.GET("/", dependency.PulseController.GetPulses)
			pulses.GET("/:id", dependency.PulseController.GetPulseByID)
		}

	}
}
