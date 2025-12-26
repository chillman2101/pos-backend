package middleware

import (
	"log"
	"pos-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				response.InternalServerError(c, "Internal server error", err)
				c.Abort()
			}
		}()
		c.Next()
	}
}
