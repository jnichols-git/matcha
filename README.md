# cloudretic/router

[![Coverage Status](https://coveralls.io/repos/github/cloudretic/router/badge.svg)](https://coveralls.io/github/cloudretic/router)

`cloudretic/router` is an actively developed HTTP router for Go, primarily developed for CloudRETIC's API handlers but free to use by anyone under the Apache 2.0 license.

## Features

`router` supports:

- Static string routes
- Wildcard parameters
- Regex route validation
- Partial route matching
- Middleware

## Installation

`go get github.com/cloudretic/router@version`

Supported versions:

- `latest (v0.0)`

## Basic Usage

Create routes using `route.New` or `route.Declare` and a route expression:

```go
rt, err := route.New("/static")
```

Create a router using `router.New` or `router.Declare` with a router type and configuration functions:

```go
rt, err := route.New(http.MethodGet, "/static")
if err != nil { ... }
r, err := router.New(
    router.Default(),
    router.WithRoute(rt, someHandler),
    router.WithMiddleware(someMiddleware),
    ...
)
```

See `docs/` for information on implementing more advanced features.

## Maintainers

Name | Role | Pronouns | GitHub Username | Contact
---|---|---|---|---
Jake Nichols | Creator | they/them | jakenichols2719 | jnichols@cloudretic.com
