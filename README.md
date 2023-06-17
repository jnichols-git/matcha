# matcha

[![Coverage Status](https://coveralls.io/repos/github/cloudretic/matcha/badge.svg?branch=main)](https://coveralls.io/github/cloudretic/matcha?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/cloudretic/matcha)](https://goreportcard.com/report/github.com/cloudretic/matcha)
[![Discord Badge](https://img.shields.io/badge/Join%20us%20on-Discord-blue)](https://discord.gg/gCdJ6NPm)

`cloudretic/matcha` is an HTTP router designed for ease of use, power, and extensibility.

## Features

- Flexible routing--handle your API specifications with ease
- Extensible components for edge cases and integration with 3rd-party tools
- High performance that scales to larger APIs
- Comprehensive and passing test coverage, and extensive benchmarks to track performance
- Easy conversion from standard library; uses stdlib handler signatures
- Zero dependencies, zero dependency management

For a preview of what's upcoming, see our [roadmap](docs/roadmap.md).

## Installation

`go get github.com/cloudretic/matcha@v1.2.0`

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
- 1 requirement for target host on all endpoints

The MockBoards benchmarks are run alongside an *offset benchmark* that measures the performance cost of setting up scaffolding
for each request sent to calculate their final score. The values below represent performance numbers that you might expect
to see in practice. Please keep in mind that performance varies by machine--you should run benchmarks on your own
hardware to get a proper idea of how well Matcha's performance suits your needs.

### MockBoards API Spec Benchmark

Benchmark | ns/request | B/request | allocs/request
--- | --- | --- | ---
Sequential | 2226 ns/request | 1909 bytes/request | 27 allocs/request
Concurrent | 1953 ns/request | 1943 bytes/request | 29 allocs/request

### MockBoards API Routing-Only Benchmark

This is the same spec as above, but with the non-path features stripped out to give a better idea of pure routing costs.

Benchmark | ns/request | B/request | allocs/request
--- | --- | --- | ---
Sequential | 1054 ns/request | 1417 bytes/request | 12 allocs/request
Concurrent | 1353 ns/request | 1453 bytes/request | 14 allocs/request

## Maintainers

Name | Role | Pronouns | GitHub Username | Contact
---|---|---|---|---
Jake Nichols | Creator | they/them | jakenichols2719 | <jnichols@cloudretic.com>
