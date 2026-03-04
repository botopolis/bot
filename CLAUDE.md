# CLAUDE.md

## Project Overview

**botopolis** is a Hubot-inspired chatbot framework written in Go. It provides a plugin-based architecture for building bots that integrate with various chat services (e.g., Slack).

## Build & Test Commands

```bash
# Build
go build ./...

# Run all tests (with race detection)
go test -race -timeout 30s ./...

# Run tests with coverage
go test -race -timeout 30s -coverprofile=c.out ./...

# Run a specific package's tests
go test -race ./help/...
```

## Architecture

### Core Components

- **`robot.go`** — Central `Robot` struct; entry point via `bot.New(chat, plugins...)` and `robot.Run()`
- **`chat.go`** — `Chat` adapter interface and `Message` types; defines `Responder` for replying
- **`brain.go`** — `Brain` key/value store with in-memory cache and pluggable `Store` backend
- **`plugin.go`** — `Plugin` interface (`Load(*Robot)`, optional `Unload(*Robot)`) and registry
- **`responder_queue.go`** — Event dispatch: routes incoming messages to registered listeners
- **`web.go`** — Built-in HTTP server (Gorilla Mux router, default port 9090)
- **`logger.go`** — Pluggep `Logger` interface

### Listeners

```go
robot.Hear(matcher, handler)    // Any message
robot.Respond(matcher, handler) // Messages directed at bot
robot.Enter(handler)            // User joins channel
robot.Leave(handler)            // User leaves channel
robot.Topic(handler)            // Topic change
```

### Matchers

```go
bot.Regexp("^hello")   // Regular expression
bot.Contains("hi")     // Substring
bot.User("alice")      // Specific user
bot.Room("general")    // Specific room
```

### Mock Package

`mock/` contains test doubles for all core interfaces:
- `mock.NewChat()` — fake chat adapter
- `mock.NewLogger()` — fake logger
- `mock.NewStore()` — fake brain store
- `mock.NewPlugin()` — fake plugin

Use these in tests instead of real adapters.

## Conventions

- Test files follow the pattern `*_test.go` (external) and `*_internal_test.go` (white-box)
- Internal tests use the package name directly (no `_test` suffix) for access to unexported symbols
- Plugin examples live in `help/`
- Plugins receive the full `*Robot` on `Load`, giving access to `Brain`, `Router`, `Logger`, and all listener APIs

## Configuration

| Environment Variable | Default | Description |
|---|---|---|
| `PORT` | `9090` | HTTP server port |

## Dependencies

- `github.com/gorilla/mux` — HTTP routing
- `github.com/stretchr/testify` — Test assertions
