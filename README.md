# sdk-go-fiber-v2

[Retica](https://retica.sh) tracing middleware for
[gofiber/fiber v2](https://gofiber.io).

## Install

```bash
go get github.com/retica-sh/sdk-go-fiber-v2
go get github.com/retica-sh/sdk-go-retica
```

## Use

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    retica "github.com/retica-sh/sdk-go-retica"
    reticafiber "github.com/retica-sh/sdk-go-fiber-v2"
)

func main() {
    r := retica.New(
        retica.WithIngestKey("ik_live_..."),
        retica.WithServiceName("my-service"),
    )
    defer r.Shutdown()

    app := fiber.New()
    app.Use(reticafiber.Middleware(r))

    app.Get("/users/:id", func(c *fiber.Ctx) error {
        return c.SendString("ok")
    })

    app.Listen(":3000")
}
```

Every request is traced automatically. The middleware reads
`X-Retica-Trace-ID` and `X-Retica-Span-ID` from incoming requests (parent
span) and sets both on the response for downstream propagation.

## Configuration

All options live on the core `*retica.Retica`. See
[`sdk-go-retica`](https://github.com/retica-sh/sdk-go-retica) for the full
list (ingest URL, batch size, sampling, etc.) and environment variables.

## Forwarding traces to other services

```go
app.Get("/proxy", func(c *fiber.Ctx) error {
    req, _ := http.NewRequest("GET", "https://other-service/api", nil)
    req.Header.Set("X-Retica-Trace-ID", c.Get("X-Retica-Trace-ID"))
    req.Header.Set("X-Retica-Span-ID", c.Get("X-Retica-Span-ID"))
    // ... do request
})
```

## License

MIT.
