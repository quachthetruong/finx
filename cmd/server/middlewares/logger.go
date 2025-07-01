package middlewares

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"

	"financing-offer/cmd/server/request"
	"financing-offer/internal/appcontext"
	"financing-offer/pkg/timehelper"
)

func (middleware *Middleware) JsonLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.RequestURI == "/status" {
			c.Next()
			return
		}
		// Start timer
		start := time.Now()

		// Process Request
		c.Next()

		// Stop timer
		duration := timehelper.GetDurationInMilliseconds(start)

		requestId, _ := c.Get(request.RequestIdKey)
		loggerWithData := middleware.Logger.With(
			slog.String("client_ip", c.ClientIP()),
			slog.Int64("duration(ms)", duration),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.RequestURI),
			slog.Int("status", c.Writer.Status()),
			slog.String("user_name", appcontext.ContextGetUserName(c)),
			slog.String("request_id", requestId.(string)),
		)

		if c.Writer.Status() >= 500 {
			loggerWithData.Error(c.Errors.String())
		} else {
			loggerWithData.Info("")
		}
	}
}
