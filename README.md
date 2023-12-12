# matcha

[![Coverage Status](https://coveralls.io/repos/github/jnichols-git/matcha/v2/badge.svg?branch=main)](https://coveralls.io/github/jnichols-git/matcha/v2?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/jnichols-git/matcha/v2)](https://goreportcard.com/report/github.com/jnichols-git/matcha/v2)

Matcha is an HTTP router with lots of features and strong memory performance.

## Features

- Match wildcards, regex expressions, and partial routes
- Extend route validation with middleware and requirements
- Performant under load with complex APIs

## Installation

`go get github.com/jnichols-git/matcha/v2`

## Basic Usage

You can use `matcha.Router` to create a new Router and `router.HandleFunc` to handle a request path.

```go
package main

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2"
)

func sayHello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func HelloExample() {
	rt := matcha.Router()
	rt.HandleFunc(http.MethodGet, "/hello", sayHello)
	http.ListenAndServe(":3000", rt)
}
```

```sh
$ go run ./examples/ hello
$ curl localhost:3000/hello
Hello, World!
$
```

For more features, see our [User Guide](docs/user-guide.md).

## Performance

To help measure performance, we benchmark performance on a fake API spec for MockBoards.

- 18 distinct endpoints, including
  - 4 endpoints requiring authorization using a "client_id" header
  - 4 endpoints with an enumeration URI parameter (new/top posts, etc)
- 2 middleware components assigning a request ID and CORS headers
- 1 requirement for target host on all endpoints

The benchmark checks performance against single sequential requests and bursts of 10 concurrent requests. An "offset" is also calculated for the cost of building requests. The results for each benchmark are `result / count - offset`

### MockBoards

Benchmark | ns/request | B/request | allocs/request
--- | --- | --- | ---
Sequential | 2226 ns/request | 1909 bytes/request | 27 allocs/request
Concurrent | 1953 ns/request | 1943 bytes/request | 29 allocs/request

### MockBoards with v2

This mounts a copy of the API at `/v2` and runs requests against both the v1 and v2 APIs.

Benchmark | ns/request | B/request | allocs/request
--- | --- | --- | ---
Sequential | 2797 ns/request | 2046 bytes/request | 29 allocs/request
Concurrent | 2097 ns/request | 2139 bytes/request | 30 allocs/request

### MockBoards routing-only

Benchmark | ns/request | B/request | allocs/request
--- | --- | --- | ---
Sequential | 1054 ns/request | 1417 bytes/request | 12 allocs/request
Concurrent | 1353 ns/request | 1453 bytes/request | 14 allocs/request

## Maintainers

Name | Role | Pronouns | GitHub Username | Contact
---|---|---|---|---
Jake Nichols | Creator | they/them | jakenichols2719 | <mail@jnichols.info>

## License

Copyright 2023 Matcha Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

<http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
