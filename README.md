# matcha

[![Coverage Status](https://coveralls.io/repos/github/cloudretic/router/badge.svg)](https://coveralls.io/github/cloudretic/router)

`cloudretic/matcha` is an actively developed HTTP router for Go, primarily developed for CloudRETIC's API handlers but free to use by anyone under the Apache 2.0 license.

> **In version 1.1, `cloudretic/router` will be renamed to `cloudretic/matcha`.** The former name is too generic to support long-term. We apologize for any inconvenience caused by this change.

## Features

- Static string routes, wildcard parameters, regex validation, and prefix routes
- Highly customizable route/router construction; use the syntax that feels best to you
- Comprehensive and passing test coverage, and extensive benchmarks to track performance
- Native middleware to help you add common functionality, extensible when native support doesn't fit your use case
- No dependencies, what you see is what you get

## Installation

`go get github.com/cloudretic/matcha[@version]`

Supported versions:

- `v1.0`
- `main (v1.1)`
- `v1.2`

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

## Benchmarks

> These benchmarks are run on the GitHub API provided by [julienschmidt](https://github.com/julienschmidt/go-http-routing-benchmark), updated to match the current Go version.

Short answer: in tests with handling of *single requests* to a large API (~200 routes), `matcha` can handle requests end-to-end in about 470 nanoseconds, using about 720 bytes of memory, when running on an M2 MacBook Pro.

Long answer: Go benchmarks provide a measurement of `ns/op` and `B/op`, representing how much time and memory was used for one "operation", which in this case is one full loop of handling *every route* in the API, a common metric used to compare http routers in Go. Since speed in nanoseconds can be machine-dependent, we have provided a relative value instead for this comparison, where the value is (`matcha` result)/(`other` result). Higher is better/faster.

Router Name | Relative Speed | Memory Use
--- | --- | ---
[`gorilla/mux`](https://github.com/gorilla/mux) | .06x | 199,686 B/op
`matcha` | 1.0x | 139,064 B/op
[`chi`](https://github.com/go-chi/chi) | 1.52x | 61,713 B/op
[`httprouter`](https://github.com/julienschmidt/httprouter) | 5.87x | 13,792
[`gin`](https://github.com/gin-gonic/gin) | 5.87x | 0 B/op

## Maintainers

Name | Role | Pronouns | GitHub Username | Contact
---|---|---|---|---
Jake Nichols | Creator | they/them | jakenichols2719 | jnichols@cloudretic.com
