package reticafiber

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	retica "github.com/retica-sh/sdk-go-retica"
)

func setupTestApp(r *retica.Retica) *fiber.App {
	app := fiber.New()
	app.Use(Middleware(r))
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("world")
	})
	app.Get("/error", func(c *fiber.Ctx) error {
		return c.Status(500).SendString("boom")
	})
	return app
}

func TestMiddleware_SetsTraceHeaders(t *testing.T) {
	r := retica.New()
	t.Cleanup(r.Shutdown)
	app := setupTestApp(r)

	resp, err := app.Test(httptest.NewRequest("GET", "/hello", nil))
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if got := resp.Header.Get("X-Retica-Trace-ID"); len(got) != 32 {
		t.Errorf("trace ID length = %d, want 32", len(got))
	}
	if got := resp.Header.Get("X-Retica-Span-ID"); len(got) != 16 {
		t.Errorf("span ID length = %d, want 16", len(got))
	}
}

func TestMiddleware_PropagatesIncomingTrace(t *testing.T) {
	r := retica.New()
	t.Cleanup(r.Shutdown)
	app := setupTestApp(r)

	const incoming = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	req := httptest.NewRequest("GET", "/hello", nil)
	req.Header.Set("X-Retica-Trace-ID", incoming)
	req.Header.Set("X-Retica-Span-ID", "bbbbbbbbbbbbbbbb")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if got := resp.Header.Get("X-Retica-Trace-ID"); got != incoming {
		t.Errorf("trace ID = %q, want %q", got, incoming)
	}
}

func TestMiddleware_ErrorResponseStillReturnsHandlerError(t *testing.T) {
	r := retica.New()
	t.Cleanup(r.Shutdown)
	app := setupTestApp(r)

	resp, err := app.Test(httptest.NewRequest("GET", "/error", nil))
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != 500 {
		t.Errorf("status = %d, want 500", resp.StatusCode)
	}
}
