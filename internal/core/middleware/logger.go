package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		method := c.Request.Method

		// Determine log level based on status code
		var event *zerolog.Event
		if statusCode >= 500 {
			event = log.Error()
		} else if statusCode >= 400 {
			event = log.Warn()
		} else {
			event = log.Info()
		}

		// Log out the request metadata
		event.
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Str("latency", latency.String()).
			Str("ip", c.ClientIP())

		if query != "" {
			event.Str("query", query)
		}

		event.Msg("api request")
	}
}
