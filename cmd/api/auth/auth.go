package auth

import (
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		const Bearer_schema = "Bearer "
		header := context.GetHeader("Authorization")

		if header == "" {
			context.AbortWithStatus(401)
			return
		}

		token := header[len(Bearer_schema):]

		claims, err := NewJWTService().ValidateToken(token)
		if err != nil {
			context.AbortWithStatus(401)
			return
		}

		context.Set("user_id", claims.Sum)

		context.Next()
	}
}
