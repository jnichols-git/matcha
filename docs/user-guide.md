# Matcha User Guide

- [Basics](#basics)
  - [Hello World](#hello-world)
  - [Echo Server with Route Parameters](#echo-server-with-route-parameters)
  - [File Server with Partial Routes](#file-server-with-partial-routes)
  - [Note: Registration Order](#note-registration-order)
- [Advanced Usage](#advanced-usage)
  - [Mounting Subrouters](#mounting-subrouters)
  - [Complex Routes](#complex-routes)
  - [Middleware](#middleware)
  - [Requirements](#requirements)

Hello! This is a step-by-step guide to using Matcha for HTTP handling in Go.

## Basics

### Hello World

To start us off, here's a basic example. You can find project examples in the GitHub [here](https://github.com/jnichols-git/matcha/v2/tree/main/examples) if you want to experiment with them.

```go
package examples

import (
    "net/http"

    "github.com/jnichols-git/matcha/v2/pkg/router"
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

In this example, we use `router.Default` to create a router. This gives us the base router with no additional features. Then we call `rt.HandleFunc`, which handles GET requests to `/hello` with the function `sayHello`. Behind the scenes, Matcha constructs a Route with the method and path you provide, and registers the handler to that path. That means you can use some of the routing features that are provided through path syntax!

### Echo Server with Route Parameters

Wildcards match any token. You create them by including square brackets in your route expression. Here's an example echo server:

```go
package examples

import (
    "net/http"

    "github.com/jnichols-git/matcha/v2/pkg/rctx"
    "github.com/jnichols-git/matcha/v2/pkg/router"
)

func echoAdmin(w http.ResponseWriter, req *http.Request) {
    name := rctx.GetParam(req.Context(), "name")
    w.Write([]byte("Hello, admin " + name + "!"))
}

func echo(w http.ResponseWriter, req *http.Request) {
    name := rctx.GetParam(req.Context(), "name")
    w.Write([]byte("Hello, " + name + "!"))
}

func EchoExample() {
    rt := router.Default()
    rt.HandleFunc(http.MethodGet, "/hello/[name]{admin:.+}", echoAdmin)
    rt.HandleFunc(http.MethodGet, "/hello/[name]", echo)
    http.ListenAndServe(":3000", rt)
}

```

The bit in square brackets will match any value (but not *no* value) and save it in the request context under "name". You can use the `rctx` package to fetch this value. If you want to filter which values are matched, you can use regex enclosed in square brackets, like with echoAdmin.

It's important to put the echoAdmin route first here. Route are handled in the order that they are registered, and the echo route matches *anything* in the second path spot, so if they were reversed, everything would just match echo.

### File Server with Partial Routes

Routes can be configured to match their root and longer request paths by using a plus symbol `+` in the last part. This can be combined with wildcards to store the full matched path (or empty if matching the root), or regex to individually validate each path component. Here's a file server demonstrating this functionality:

```go
package examples

import (
    "net/http"
    "os"

    "github.com/jnichols-git/matcha/v2/pkg/rctx"
    "github.com/jnichols-git/matcha/v2/pkg/router"
)

type fileServer struct {
    root string
}

func (fs *fileServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    path := rctx.GetParam(req.Context(), "filepath")
    dat, err := os.ReadFile(fs.root + path)
    if err != nil {
        w.WriteHeader(404)
        w.Write([]byte("File " + path + " does not exist."))
        return
    }
    w.Write(dat)
}

func FileServer(dir string) {
    rt := router.Default()
    rt.Handle(http.MethodGet, "/files/[filepath]+", &fileServer{dir})
    http.ListenAndServe(":3000", rt)
}

```

Since ***routes are matched in the order they are registered***, partial routes will override any longer path you register to a router afterwards. This is particularly important for partial routes. You should register these last, and if you have multiple, in order from longest to shortest.

### Note: Registration Order

Given the emphasis put onto registration order here, I think it's important to note *why* Matcha works this way. When you register a route, Matcha adds it to a tree made up of the parts between the slashes. This tree is traversed depth-first and in order, and the first match is returned immediately, meaning that only some subset of the routes you register are checked on any incoming request. This is very fast.

Implicitly deprioritizing some routes to skew towards exact matches causes two problems with this structure:

- The ordering of routes is no longer guaranteed; Matcha shouldn't be doing anything you don't explicitly tell it to do
- Performance will be decreased as in order to know if the most exact match has been reached, the entire tree must be traversed

We're working on ways to make routing more intuitive while avoiding these problems. In the meantime, we believe that strict registration order is the best way to go, so that you can always predict what Matcha will do with the instructions you give it.

## Advanced Usage

### Mounting Subrouters

In addition to `router.Handle` and similar functions, there is also `router.Mount`, which passes all requests to a path prefix through to the given handler. In this example, any request made to `localhost:3000/api` will go directly to the `api2` router.

```go
func main() {
    api1 := router.Default()
    api2 := router.Default()
    // Register some handlers here.
    api1.Mount("/v2", api2)
    http.ListenAndServe(":3000", api1)
}
```

Usage notes:

- This automatically trims the given prefix, so routes in `api2` do *not* need to include `/v2` in their paths.
- Mounting currently only supports static string paths.
- The underlying path used for mounting is a partial path, and comes with all of the same caveats.

### Complex Routes

So, what if you need more out of your routes?

Behind the scenes, `Handle` and `HandleFunc` use the method and path you provide do construct a route and register the handler to it. This covers a lot of use cases, but some applications need more control over the behavior of a route. For this, we provide `route.New` and `route.Declare`, which both accept a variadic list of arguments modifying the route. These are called `ConfigFunc`s, and they give access to things like middleware or "requirements", which match against non-path properties of a request. `HandleRoute` and `HandleRouteFunc` are used to register these routes directly.

Our [routes documentation](routes.md) has the most up-to-date information on advanced usage, but here are some commonly-used examples.

### Middleware

Matcha uses `func(w http.ResponseWriter, req *http.Request) *http.Request` for middleware. You can attach them to a router or route using `Attach` or the ConfigFunc `WithMiddleware` respectively. This example logs all incoming requests, and rejects requests to the second route that don't have a query parameter `user` by returning `400 Bad Request`:

```go
package main

import (
    "github.com/jnichols-git/matcha/v2/pkg/router"
    "github.com/jnichols-git/matcha/v2/pkg/route"
    "github.com/jnichols-git/matcha/v2/pkg/middleware"
)

func main() {
    userRoute := route.Declare(
        http.MethodGet, "/users",
        route.WithMiddleware(middleware.ExpectQueryParam("user"))
    )
    rt := router.Default()
    rt.Handle(http.MethodGet, "/", h1),
    rt.HandleRoute(userRoute, h2)
    server.Attach(middleware.LogRequests())
}
```

### Requirements

Matcha provides an interface for matching things that are not paths in package `route/require`. You can define your own with the function definition `func(req *http.Request) bool` and register them onto routes by using the config function or route function `route.Require`. If a requirement returns `false`, the router will continue to match against the remaining routes.

```go
webRoute, err := route.New(
    http.MethodGet, "/",
    route.Require(require.HostPorts("https://{www.|}decentplatforms.com")),
)
apiRoute, err := route.New(
    http.MethodGet, "/",
    require.HostPorts("https://api.decentplatforms.com"),
)
```
