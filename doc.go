// Package reticafiber is the Retica SDK adapter for gofiber/fiber v2.
//
// # Quick start
//
//	import (
//	    "github.com/gofiber/fiber/v2"
//	    retica "github.com/retica-sh/sdk-go-retica"
//	    reticafiber "github.com/retica-sh/sdk-go-fiber-v2"
//	)
//
//	func main() {
//	    app := fiber.New()
//
//	    r := retica.New(
//	        retica.WithIngestKey("ik_live_..."),
//	        retica.WithServiceName("my-service"),
//	    )
//	    defer r.Shutdown()
//
//	    app.Use(reticafiber.Middleware(r))
//	    app.Listen(":3000")
//	}
//
// All configuration (ingest key, service name, batching, sampling) is on the
// core *retica.Retica — see github.com/retica-sh/sdk-go-retica for options.
//
// # Trace propagation
//
// The middleware reads X-Retica-Trace-ID and X-Retica-Span-ID from incoming
// requests (treating the latter as the parent span) and sets both headers on
// the response so downstream services can inherit the trace. When forwarding
// requests to other services, copy the headers from the current fiber.Ctx.
package reticafiber
