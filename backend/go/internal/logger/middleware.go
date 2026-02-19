package logger

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Middleware returns a Fiber middleware that injects a request ID and logs
// basic request/response information using the provided logger.
func Middleware(log Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Ensure a request ID is present.
		reqID := c.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
			c.Set("X-Request-ID", reqID)
		}

		// Add request ID to context and logger.
		ctx := context.WithValue(c.UserContext(), ContextKeyRequestID, reqID)
		c.SetUserContext(ctx)

		// Capture request body
		reqBody := string(c.Body())
		
		l := log.With(
			"request_id", reqID,
			"path", c.Path(),
			"method", c.Method(),
			"ip", c.IP(),
			"user_agent", c.Get("User-Agent"),
			"request_headers", c.GetReqHeaders(),
			"query_params", c.Queries(),
			"request_body", reqBody,
		)

		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()
		respBody := string(c.Response().Body())

		fields := []interface{}{
			"status", status,
			"duration_ms", duration.Milliseconds(),
			"response_headers", c.GetRespHeaders(),
			"response_body", respBody,
		}

		if err != nil {
			fields = append(fields, "error", err)
			l.Error("request completed with error", fields...)
			return err
		}

		l.Info("request completed", fields...)
		return nil
	}
}


