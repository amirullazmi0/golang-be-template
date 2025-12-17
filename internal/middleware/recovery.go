package middleware

import (
	"net/http"

	"github.com/amirullazmi0/kratify-backend/pkg/logger"
	"github.com/amirullazmi0/kratify-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery middleware recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("request_id", c.GetString("RequestID")),
				)

				response.Error(c, http.StatusInternalServerError, "Internal server error", nil)
				c.Abort()
			}
		}()

		c.Next()
	}
}
