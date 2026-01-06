package middleware

import (
	"time"

	"github.com/amirullazmi0/kratify-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger middleware logs HTTP requests with structured fields for Grafana Loki
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		// Get error if exists
		var errorMsg string
		if len(c.Errors) > 0 {
			errorMsg = c.Errors.String()
		}

		// Structured logging for Grafana Loki
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", c.Writer.Status()),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.Int64("latency_ms", latency.Milliseconds()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("request_id", c.GetString("RequestID")),
			zap.Int("response_size", c.Writer.Size()),
		}

		if errorMsg != "" {
			fields = append(fields, zap.String("error", errorMsg))
		}

		// Log with appropriate level based on status code
		status := c.Writer.Status()
		switch {
		case status >= 500:
			logger.Error("HTTP Request - Server Error", fields...)
		case status >= 400:
			logger.Warn("HTTP Request - Client Error", fields...)
		default:
			logger.Info("HTTP Request", fields...)
		}
	}
}
