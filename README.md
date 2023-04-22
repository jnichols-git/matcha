# matcha

[![Coverage Status](https://coveralls.io/repos/github/cloudretic/router/badge.svg)](https://coveralls.io/github/cloudretic/router)

`cloudretic/matcha` is an actively developed HTTP router for Go, primarily developed for CloudRETIC's API handlers but free to use by anyone under the Apache 2.0 license.

## Features

- Static string routes, wildcard parameters, regex validation, and prefix routes
- Highly customizable route/router construction; use the syntax that feels best to you
- Comprehensive and passing test coverage, and extensive benchmarks to track performance
- Native middleware
  - Log inbound requests
  - Require query parameters
  - Attach middleware to routes to target specific endpoints
- CORS handling
  - Native middleware to write CORS headers on responses
  - A premade handler to manage OPTIONS requests

## Benchmarks

These benchmarks are run on the GitHub API provided by [julienschmidt](https://github.com/julienschmidt/go-http-routing-benchmark), updated to match the current Go version.

Go benchmarks provide a measurement of `ns/op` and `B/op`, representing how much time and memory was used for one "operation", which in this case is one full loop of handling *every route* in the API. Since speed in nanoseconds can be machine-dependent, we have provided a relative value instead, where the value is (`matcha` result)/(`other` result). Higher is better/faster.

Router Name | Relative Speed | Memory Use
--- | --- | ---
[`gorilla/mux`](https://github.com/gorilla/mux) | .06x | 199,686 B/op
`matcha` | 1.0x | 139,064 B/op
[`chi`](https://github.com/go-chi/chi) | 1.52x | 61,713 B/op
[`httprouter`](https://github.com/julienschmidt/httprouter) | 5.87x | 13,792
[`gin`](https://github.com/gin-gonic/gin) | 5.87x | 0 B/op

## Installation

`go get github.com/cloudretic/matcha[@version]`

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
