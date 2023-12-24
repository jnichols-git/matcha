# Matcha User Guide

- [Basics](#basics)
	- [Hello World](#hello-world)
	- [Echo Server with Route Parameters](#echo-server-with-route-parameters)
	- [File Server with Partial Routes](#file-server-with-partial-routes)
- [Advanced Usage](#advanced-usage)
	- [Mounting Subrouters](#mounting-subrouters)
	- [Middleware](#middleware)
	- [Requirements](#requirements)
- [FAQ](#faq)
	- [Why are some of my routes not matched when they should be?](#why-are-some-of-my-routes-not-matched-when-they-should-be)

Hello! This is a step-by-step guide to using Matcha for HTTP handling in Go.

## Basics

There are a few examples in the `examples` directory to show basic usage.

### Hello World

You can use `matcha.Router` to create a new Router and `router.HandleFunc` to
handle a request path. Routers are type `*router.Router` if you need to store
it.

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

### Echo Server with Route Parameters

If a path segment starts with `:`, it will match *any* value, and incoming requests can access the matched value using `matcha.RouteParam`.

```go
package main

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2"
)

func echo(w http.ResponseWriter, req *http.Request) {
	name := matcha.RouteParam(req, "name")
	w.Write([]byte("Hello, " + name + "!"))
}

func EchoExample() {
	rt := matcha.Router()
	rt.HandleFunc(http.MethodGet, "/hello/:name", echo)
	http.ListenAndServe(":3000", rt)
}
```

```sh
$ go run ./examples/ echo
$ curl localhost:3000/hello/jnichols
Hello, jnichols!
$
```

### File Server with Partial Routes

If a route ends in `/+`, it will match its root value + any "tail" values. You can also do `/:param+` to store the "tail" in a RouteParam.

```go
package main

import (
	"net/http"
	"os"

	"github.com/jnichols-git/matcha/v2"
)

type fileServer struct {
	root string
}

func (fs *fileServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := matcha.RouteParam(req, "filepath")
	dat, err := os.ReadFile(fs.root + path)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("File " + path + " does not exist."))
		return
	}
	w.Write(dat)
}

func FileServer(dir string) {
	rt := matcha.Router()
	rt.Handle(http.MethodGet, "/files/:filepath+", &fileServer{dir})
	http.ListenAndServe(":3000", rt)
}
```

```sh
$ go run ./examples/ fileserver
$ curl localhost:3000/files/data/hello.txt
Hello, fileserver!
$
```

## Advanced Usage

### Mounting Subrouters

```go
func main() {
    api1 := matcha.Router()
    api2 := matcha.Router()
    // Register some handlers here.
    api1.Mount("/v2", api2)
    http.ListenAndServe(":3000", api1)
}
```

You can use `router.Mount` to pass *all* requests to a handler. `/v2` is treated as a prefix; it is matched as a partial route and removed from requests before passing to the mounted handler.

### Middleware

Middleware is defined in `teaware` as a `func(next http.Handler) http.Handler`. You can `Use` middleware with any Route or Router to perform actions before handlers. Routers will run middleware before routing a request, and Routes will run it before their handlers. Returning `nil` will reject the request.

```go
package main

import (
	"net/http"

	"github.com/jnichols-git/matcha/v2"
)

func ValidateName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if name := matcha.RouteParam(r, "name"); name[0] != 'A' {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Names must start with 'A'.\n"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func MiddlewareExample() {
	router := matcha.Router()
	nameRoute, _ := matcha.Route(http.MethodGet, "/hello/:name")
	nameRoute.Use(ValidateName)
	router.HandleRouteFunc(nameRoute, echo)
	http.ListenAndServe(":3000", router)
}

```

```sh
$ go run ./examples/ middleware
$ curl -v localhost:3000/hello/jnichols
< HTTP/1.1 400 Bad Request
Names must start with 'A'.
$ curl -v localhost:3000/hello/Alex
< HTTP/1.1 200 OK
Hello, Alex!
$
```

### Requirements

Requirements are defined in `require` as a `func(*http.Request) bool`. If you need more tools to check routes, you can `Require` a requirement with any Route. If this function returns false for any request, the route is not matched, but it can still match any routes after.

```go
webRoute, err := matcha.Route(
    http.MethodGet, "/",
).Require(require.HostPorts("https://[www.|]jnichols.info"))
apiRoute, err := matcha.Route(
    http.MethodGet, "/",
).Require(require.HostPorts("https://api.jnichols.info"))
```

## FAQ

### Why are some of my routes not matched when they should be?

Many routers check routes in order of *specificity*; more general routes, like wildcards, are checked after normal ones. Matcha matches in the order the routes are given to the router, so wildcards or partial routes can sometimes override normal routes.
