package reticafiber

import (
	"github.com/gofiber/fiber/v2"
	retica "github.com/retica-sh/sdk-go-retica"
)

// Middleware wraps a *retica.Retica as a fiber middleware handler. Every
// request is traced via retica.TraceHTTP and the resulting trace/span IDs
// are propagated as X-Retica-Trace-ID and X-Retica-Span-ID response headers.
func Middleware(r *retica.Retica) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := retica.RequestInfo{
			Method:       c.Method(),
			Path:         c.Path(),
			ClientIP:     c.IP(),
			UserAgent:    c.Get("User-Agent"),
			TraceID:      c.Get("X-Retica-Trace-ID"),
			ParentSpanID: c.Get("X-Retica-Span-ID"),
			RequestSize:  len(c.Request().Body()),
		}

		var runErr error
		traceID, spanID := r.TraceHTTP(req, func() retica.ResponseInfo {
			runErr = c.Next()
			return retica.ResponseInfo{
				StatusCode:   c.Response().StatusCode(),
				ResponseSize: len(c.Response().Body()),
				Err:          runErr,
			}
		})

		if traceID != "" {
			c.Set("X-Retica-Trace-ID", traceID)
			c.Set("X-Retica-Span-ID", spanID)
		}
		return runErr
	}
}
