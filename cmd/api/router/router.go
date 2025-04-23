package router

import (
	"github.com/Edu4rdoNeves/ingestor-magalu/internal/dependency"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) *gin.Engine {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600,
	}))

	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Status(204)
	})

	main := router.Group("api/v1")
	{
		pulse := main.Group("pulse")
		{
			pulse.GET("/", dependency.PulseDependency.GetPulses)
			pulse.GET("/:id", nil)
		}
	}
	return router
}
