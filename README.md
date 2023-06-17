# matcha

[![Coverage Status](https://coveralls.io/repos/github/cloudretic/matcha/badge.svg?branch=main)](https://coveralls.io/github/cloudretic/matcha?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/cloudretic/matcha)](https://goreportcard.com/report/github.com/cloudretic/matcha)
[![Discord Badge](https://img.shields.io/badge/Join%20us%20on-Discord-blue)](https://discord.gg/gCdJ6NPm)

`cloudretic/matcha` is an actively developed HTTP router for Go with a focus on providing a flexible and performant route API.

## Features

- Static string routes, wildcard parameters, regex validation, and prefix routes
- Highly customizable route/router construction; get the results you want with the syntax that feels best to you
- Comprehensive and passing test coverage, and extensive benchmarks to track performance
- Native middleware to help you add common functionality, extensible when native support doesn't fit your use case
- No dependencies, what you see is what you get

For a preview of what's upcoming, see our [roadmap](docs/roadmap.md).

## Installation

`go get github.com/cloudretic/matcha@v1.1.2`

## Basic Usage

Here's a "Hello, World" example to introduce you to Matcha's syntax! It serves requests to `http://localhost:8080/hello`

```go
package examples

import (
    "net/http"

    "github.com/cloudretic/matcha/pkg/router"
)

func sayHello(w http.ResponseWriter, req *http.Request) {
    w.Write([]byte("Hello, World!"))
}

func HelloExample() {
    rt := router.Default()
    rt.HandleFunc(http.MethodGet, "/hello", sayHello)
    // or:
    // rt.Handle(http.MethodGet, "/hello", http.HandlerFunc(sayHello))
    http.ListenAndServe(":3000", rt)
}
```

For a step-by-step guide through Matcha's features, see our [User Guide](docs/user-guide.md).

## Performance

Matcha has an extensive benchmark suite to help identify, document, and improve performance over time. Additionally, `/bench` contains a comprehensive benchmark API for "MockBoards", a fake website that just so happens to use all of the features of Matcha. The MockBoards API has the following:

- 18 distinct endpoints, including
  - 4 endpoints requiring authorization using a "client_id" header
  - 4 endpoints with an enumeration URI parameter (new/top posts, etc)
- 2 middleware components assigning a request ID and CORS headers
- 1 requirement for target host

Our benchmark constructs and then runs single requests against the MockBoards API specification, first in sequence and then in parallel. The below is the results on an M2 MacBook Pro, provided for convenience; we encourage you to test relevant benchmarks on your own hardware if you're comparing to other solutions. *Please note that 1 op is 10 requests for the concurrent benchmark.*

Benchmark | ns/op | B/op | allocs/op
--- | --- | --- | ---
Sequential | 2461 ns/request | 6032 bytes/request | 29 allocs/request
Concurrent | 2854 ns/request | 7107 bytes/request | 40 allocs/request

## Maintainers

Name | Role | Pronouns | GitHub Username | Contact
---|---|---|---|---
Jake Nichols | Creator | they/them | jakenichols2719 | <jnichols@cloudretic.com>
