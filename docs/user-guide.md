# Matcha User Guide

- [Matcha User Guide](#matcha-user-guide)
  - [Basics](#basics)
    - [Creating a Route](#creating-a-route)
    - [Creating a Router](#creating-a-router)
    - [ConfigFuncs and Declare](#configfuncs-and-declare)
  - [Route Syntax](#route-syntax)
    - [Wildcard Parameters](#wildcard-parameters)
    - [Regex Validation](#regex-validation)
    - [Partial/Prefix Routes](#partialprefix-routes)
    - [Full Example](#full-example)
    - [Note: Registration Order](#note-registration-order)
  - [Additional Tools](#additional-tools)
    - [Middleware](#middleware)

Hello! This is a step-by-step guide to using Matcha for HTTP handling in Go.

## Basics

### Creating a Route

You create a route by using `route.New` or `route.Declare`. Both create routes in the same way, but only `New` returns an error if route creation fails; Declare will panic instead. Route creation takes in a *method* and a *route expression*, and will fail if the expression is invalid. Notably, routes don't take an `http.Handler`. This happens later, so don't worry about that just yet. Here are a few basic routes:

```go
r1, err := route.New(http.MethodGet, "/")
r2 := route.Declare(http.MethodGet, "/some/route")
```

### Creating a Router

Just like with routes, you create a router using `router.New` or `route.Declare`. These take in a router object and output... a router object.

```go
rt1, err := router.New(router.Default())
rt2 := router.Declare(router.Default())
```

You can register the routes you created using `AddRoute` and a handler. **Routes will be matched in the order you add them**. Assume we've got a couple of `http.HandlerFunc`s called h1 and h2 for brevity.

```go
rt1.AddRoute(r1, h1)
rt1.AddRoute(r2, h2)
```

Now `rt1` will serve requests to the root URL and to `/some/route` using h1 and h2. Here's the full example:

```go
r1, err := r.New(http.MethodGet, "/")
if err != nil {
    ...
}
r2, err := r.New(http.MethodGet, "/some/route")
if err != nil {
    ...
}
rt1, err := router.New(router.Default())
if err != nil {
    ...
}
rt1.AddRoute(r1, h1)
rt1.AddRoute(r2, h2)
```

You may have noticed that this is a lot of error checking that you don't want to be doing. That is fair! This error behavior is meant to help you deal with complex workloads, like loading a configuration file with lots of paths to register, where errors may occur and be recoverable. Let's look at a different way to do this.

### ConfigFuncs and Declare

You may be wondering what the deal is with `router.New`. In this example, it serves only to pass through `router.Default()` and make your code longer--why go through the trouble? You certainly don't have to, but in addition to their regular arguments, both `route.New` and `router.New` take a variadic slice of functions that run on their respective structures. These are used to further customize behavior and enable a *declarative* routing style, which is what `Declare` is for. Here's the above example, but done in this style:

```go
rt1 := router.Declare(
    router.Default(),
    router.WithRoute(r.Declare(http.MethodGet,"/"), h1),
    router.WithRoute(r.Declare(http.MethodGet,"/some/route"), h2),
)
```

## Route Syntax

Now, let's talk about *route syntax*. In the examples above, our routes are made up entirely of *static parts*, which means that every token contained between slashes `/` or the end of the route is a URL-encoded string. These will match exactly with incoming requests. However, there are some additional features you can use to customize how routes behave!

### Wildcard Parameters

Wildcards match any token. You create them by including square brackets in your route expression, like this:

```go
r1 := route.Declare(http.MethodGet, "/files/[filename]")
r2 := route.Declare(http.MethodGet, "/users/[id]")
```

The value that is matched by the wildcard is stored for your use later, and you can access them by using `rctx.GetParam(r.Context, "paramName")` in your handler. Since parameters don't match an empty string, these are guaranteed to contain values if the route is matched. Here's an example for the route `r1` above.

```go
func h1(w http.RequestWriter, req *http.Request) {
    fn := rctx.GetParam(req.Context, "filename")
}
```

Since **routes are matched in the order they are registered**, wildcards will override any same-length path you register to a router afterwards.

### Regex Validation

If you have a particular part of the route you want to ensure follows a specific format, you can use regex to reject any non-matchingr request. Any pattern contained in squiggly brackets `{}` will be handled as regex. You can even combine this with a wildcard to create an auto-validated parameter! If you do this, the entire token will be matched--groups aren't taken into account for parameters.

```go
r1 := r.Declare(http.MethodGet, `/{hello|goodbye}`)
r2 := r.Declare(http.MethodGet, `/files/[filename]{.*\.(md|go)}`)
```

Since **routes are matched in the order they are registered**, permissive regex will override any same-length path you register to a router afterwards.

### Partial/Prefix Routes

Routes can be configured to match their root and longer request paths by using a plus symbol `+` in the last part. This can be combined with wildcards to store the full matched path (or empty if matching the root), or regex to individually validate each path component.

```go
r1 := r.Declare(http.MethodGet, `/files/[filename]+`)
r2 := r.Declare(http.MethodGet, `/+`)
```

Since--and I promise this is the last time we'll say this--***routes are matched in the order they are registered***, partial routes will override any longer path you register to a router afterwards. This is particularly important for partial routes. You should register these last, and if you have multiple, in order from longest to shortest.

### Full Example

In this example, we have 4 handlers:

- `indexHandler`, which returns an HTML file for a website homepage
- `reviewCreate`, which allows the user to POST a review with a string name
- `reviewGet`, which GETs a review
- `staticHandler`, which serves static files

```go
import (
    rt "github.com/cloudretic/matcha/pkg/router"
    r "github.com/cloudretic/matcha/pkg/route"
)

/* handlers defined here */

func main() {
    server := rt.Declare(
        rt.Default(),
        rt.WithRoute(r.Declare(http.MethodGet, "/"), indexHandler),
        rt.WithRoute(r.Declare(http.MethodPost, "/reviews/[name]"), reviewCreate),
        rt.WithRoute(r.Declare(http.MethodGet, "/reviews/[name]"), reviewGet),
        rt.WithRoute(r.Declare(http.MethodGet, `/static/[filename]{\w+.(.*)?}+`), staticHandler),
    )
    http.ListenAndServe(":3000", server)
}
```

### Note: Registration Order

Given the emphasis put onto registration order here, I think it's important to note *why* Matcha works this way. When you register a route, Matcha adds it to a tree made up of the parts between the slashes. This tree is traversed depth-first and in order, and the first match is returned immediately, meaning that only some subset of the routes you register are checked on any incoming request. This is very fast.

Implicitly deprioritizing some routes to skew towards exact matches causes two problems with this structure:

- The ordering of routes is no longer guaranteed; Matcha shouldn't be doing anything you don't explicitly tell it to do
- Performance will be hit as in order to know if the most exact match has been reached, the entire tree must be traversed

We're working on ways to make routing more intuitive while avoiding these problems. In the meantime, we believe that strict registration order is the best way to go, so that you can always predict what Matcha will do with the instructions you give it.

## Additional Tools

### Middleware

Matcha uses `func(w http.ResponseWriter, req *http.Request) *http.Request` for middleware. You can attach them to a router or route using `Attach`, or the ConfigFunc `WithMiddleware`. This example logs all incoming requests, and rejects requests to the second route that don't have a query parameter `user` by returning `400 Bad Request`:

```go
import (
    "github.com/cloudretic/matcha/pkg/router"
    "github.com/cloudretic/matcha/pkg/route"
    "github.com/cloudretic/matcha/pkg/middleware"
)

server := router.Declare(
    router.Default(),
    router.WithRoute(route.Declare(http.MethodGet, "/"), h1),
    router.WithRoute(route.Declare(
        http.MethodGet, "/users",
        WithMiddleware(middleware.ExpectQueryParam("user"))
    ), h2)
)
server.Attach(middleware.LogRequests())
```

Check `package middleware` for information on what we natively support. Additionally, in version 1.2.0, we are extending support to `http.Handler` middleware to more easily integrate with exterior tools.
