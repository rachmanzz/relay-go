# Relay: SSE Driver & Long Polling Module

Relay is a robust Go communication library that provides reliable real-time data streaming by combining Server-Sent Events (SSE) with a stateful Long Polling fallback.

## Features

- **Hybrid Transport**: SSE primary with automatic stateful Long Polling fallback.
- **Framework Agnostic**: Works seamlessly with Goravel, Fiber, Echo, Gin, and standard `net/http`.
- **Hybrid Storage**: Support for In-Memory (default) and Redis session management.
- **Intelligent Driver**: Go client that manages transport switching automatically.

## Installation

```bash
go get github.com/rachmanzz/relay-go
```

## Quick Start (Server)

```go
import "github.com/rachmanzz/relay-go"

// Initialize Relay (In-Memory)
r := relay.New()

// Initialize Relay (Redis)
// r := relay.New(relay.WithRedis(redisClient))
```

### Integration Examples

#### Goravel Integration
```go
func (c *NotificationController) GetNotifications(ctx http.Context) http.Response {
    pollingID := ctx.Request().Input("polling_id")
    
    result := r.Polling(ctx.Context(), pollingID, 40, func(lastTS time.Time) any {
        // Your logic to check for new data since lastTS
        return data
    })
    
    return ctx.Response().Json(200, result)
}
```

#### Fiber Integration
```go
app.Get("/polling", func(c *fiber.Ctx) error {
    pollingID := c.Query("polling_id")
    
    result := r.Polling(c.Context(), pollingID, 30, func(ts time.Time) any {
        return getData(ts)
    })
    
    return c.JSON(result)
})
```

#### SSE Handler (Standard http.Handler)
```go
// Most frameworks allow attaching a standard http.Handler
http.Handle("/events", r.SSEHandler())
```

## Using the Intelligent Driver (Client)

```go
driver := relay.NewDriver("http://server/events", "http://server/polling")
driver.Start(context.Background())

for event := range driver.Events() {
    fmt.Printf("Received: %s - %v\n", event.Name, event.Data)
}
```

## Configuration Options

- `WithTTL(time.Duration)`: Set session expiration time (default 5m).
- `WithRedis(*redis.Client)`: Use Redis for distributed session management.
