# cloudretic/router

[![Coverage Status](https://coveralls.io/repos/github/cloudretic/router/badge.svg)](https://coveralls.io/github/cloudretic/router)

`cloudretic/router` is an actively developed HTTP router for Go, primarily developed for CloudRETIC's API handlers but free to use by anyone under the Apache 2.0 license.

> **In version 1.1, `cloudretic/router` will be renamed to `cloudretic/matcha`.** The former name is too generic to support long-term. We apologize for any inconvenience caused by this change.

## Features

- Static string routes, wildcard parameters, regex validation, and prefix routes
- Highly customizable route/router construction; use the syntax that feels best to you
- Comprehensive and passing test coverage, and extensive benchmarks to track performance

### Upcoming in v1.1

- More native middleware
  - Log inbound requests
  - Require query parameters
  - Attach middleware to routes to target specific endpoints
- CORS handling
  - Native middleware to write CORS headers on responses
  - A premade handler to manage OPTIONS requests
- Significant performance improvements
  - GitHub API benchmark performs 2.5x faster than v1.0

## Installation

`go get github.com/cloudretic/router[@version]`

Supported versions:

- `main (v1.0)`
- `v1.1`

## Basic Usage

Create routes using `route.New` or `route.Declare` and a route expression:

```go
rt, err := route.New(http.MethodGet, "/static")
if err != nil { ... }
```

Create a router using `router.New` or `router.Declare` with a router type and configuration functions:

```go
// Traditional router construction; create router and add properties
r := router.Default()
rt, err := route.New(http.MethodGet, "/static")
if err != nil { ... }
r.AddRoute(rt, routeHandler)
r.AddNotFound(notFoundHandler)
r.Attach(someMiddleware)
```

```go
// Declarative router construction; create router by definition and panic on failure
r := router.Declare(
    router.Default(),
    router.WithRoute(route.Declare(http.MethodGet, "/static"), routeHandler),
    router.WithNotFound(notFoundHandler),
    router.WithMiddleware(someMiddleware)
)
```

See `docs/` for information on implementing more advanced features.

## Maintainers

Name | Role | Pronouns | GitHub Username | Contact
---|---|---|---|---
Jake Nichols | Creator | they/them | jakenichols2719 | jnichols@cloudretic.com
