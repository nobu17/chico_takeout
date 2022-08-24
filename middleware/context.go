package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
)

func SetContext(handler func(ctx context.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c.Request.Context())
		c.Next()
	}
}
