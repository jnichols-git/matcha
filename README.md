# matcha

[![Coverage Status](https://coveralls.io/repos/github/cloudretic/matcha/badge.svg)](https://coveralls.io/github/cloudretic/router)
[![Discord Badge](https://img.shields.io/badge/Join%20us%20on-Discord-blue)](https://discord.gg/gCdJ6NPm)

`cloudretic/matcha` is an actively developed HTTP router for Go, primarily developed for CloudRETIC's API handlers but free to use by anyone under the Apache 2.0 license.

## Features

- Static string routes, wildcard parameters, regex validation, and prefix routes
- Highly customizable route/router construction; use the syntax that feels best to you
- Comprehensive and passing test coverage, and extensive benchmarks to track performance
- Native middleware to help you add common functionality, extensible when native support doesn't fit your use case
- No dependencies, what you see is what you get

For a preview of what's upcoming, see our [roadmap](docs/roadmap.md).

## Installation

`go get github.com/cloudretic/matcha[@version]`

Stable versions:

- `v1.0`
- `v1.1.0 (stable release, recommended)`

Omitting the version will fetch the main branch, which may contain unreleased but stable features.

## Basic Usage

Here's a "Hello, World" example to introduce you to Matcha's syntax! It serves requests to `http://localhost:8080/hello`

```go
package main

import (
    "github.com/cloudretic/matcha/pkg/route"
    "github.com/cloudretic/matcha/pkg/router"
)

func handleHello(w http.ResponseWriter, req *http.Request) {
    w.Write([]byte("Hello, World"))
}

func main() {
    helloRoute := route.Declare(http.MethodGet, "/hello")
    s := router.Declare(
        router.Default(),
        router.WithRoute(helloRoute, http.HandleFunc(handleHello)),
    )
    http.ListenAndServer(":8080", s)
}
```

For a step-by-step guide through Matcha's features, see our [User Guide](docs/user-guide.md).

## Benchmarks

> These benchmarks are run on the GitHub API provided by [julienschmidt](https://github.com/julienschmidt/go-http-routing-benchmark), updated to match the current Go version.

Short answer: in tests with handling of *single requests* to a large API (~200 routes), Matcha can handle requests end-to-end in about 470 nanoseconds, using about 720 bytes of memory, when running on an M2 MacBook Pro.

Long answer: Go benchmarks provide a measurement of `ns/op` and `B/op`, representing how much time and memory was used for one "operation", which in this case is one full loop of handling *every route* in the API, a common metric used to compare http routers in Go. Since speed in nanoseconds can be machine-dependent, we have provided a relative value instead for this comparison, where the value is (Matcha result)/(`other` result). Higher is better/faster.

Router Name | Relative Speed | Memory Use
--- | --- | ---
[`gorilla/mux`](https://github.com/gorilla/mux) | .06x | 199,686 B/op
`cloudretic/matcha` | 1.0x | 67,928 B/op
[`go-chi/chi`](https://github.com/go-chi/chi) | 1.52x | 61,713 B/op
[`julienschmidt/httprouter`](https://github.com/julienschmidt/httprouter) | 5.87x | 13,792 B/op
[`gin-gonic/gin`](https://github.com/gin-gonic/gin) | 5.87x | 0 B/op

## Maintainers

Name | Role | Pronouns | GitHub Username | Contact
---|---|---|---|---
Jake Nichols | Creator | they/them | jakenichols2719 | <jnichols@cloudretic.com>
