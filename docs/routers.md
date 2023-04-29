# Routers

## Creating Routers

You can create a Router with a series of Routes, add Routes after creation, or both, using `WithRoute` and `AddRoute`.

```go
r := router.Declare(
    router.Default(),
    router.WithRoute(route.Declare(http.MethodGet, "/someEndpoint"), someHandler)
)
r.AddRoute(route.Declare(http.MethodGet, "/someOtherEndpoint"), someOtherHandler)
```

Routes will be handled in the order they are added. Re-adding a route with the same method and path does not change the handler or middleware for that route; only the first instance of a duplicate route will be matched against.

You can also add a specific handler that's called in the event that no route matches using `WithNotFound` or `AddNotFound` in the same way (the router will return an empty 404 by default). Doing so will override any previously-set handler for this case.

```go
r := router.Declare(
    router.Default(),
    router.WithRoute(...),
    router.WithNotFound(notFoundHandler)
)
r.AddNotFound(otherNotFoundHandler)
```

## Configuration

`ConfigFunc`s are the generalization of router customization--both `WithRoute` and `WithNotFound` are `ConfigFuncs`, and router creation functions can accept any number of them. These will be run *in order*.

## Middleware

Middleware is defined as a `func(http.ResponseWriter, *http.Request) *http.Request`. A nil return value is treated as a rejection of the request; a non-nil value will continue to be handled, and may have altered the request in some way. Common middleware options that are developed will be added as `ConfigFunc`s, but you can define your own; any function that fits that definition can serve as middleware, and can be attached using `WithMiddleware` or `Attach`:

```go
router, err := router.New(
    router.Default(),
    router.WithRoute(...),
    router.WithMiddleware(someMiddleware),
)
if err != nil { ... }
router.Attach(someOtherMiddleware)
```

## Matching

As of version 1.1, Matcha uses a tree-based system to match requests. When adding routes, the `Part`s of the route are added to a tree, where equivalent `Part`s are combined into the same node. Requested routes are compared to this tree using a depth-first traversal, in the same order that the routes were added. For example, take this tree for the following set of routes. Note that leaf nodes are not overwritten by routes that extend beyond them.

```txt
1: /static
2: /static/route/a
3: /static/route/b
4: /static/other/route

root -- static
      \ static -- route -- a
               |         \ b
                \ other -- route
```

Some routers support exact routes only, as defining behavior for conflicting routes can be expensive. Matcha supports wildcards, regex validation, and partial routes, all of which can lead to conflicts. To avoid incurring performance costs, Matcha currently supports only *first match*, which means that the first "correct" route, in order of their registration, will be matched, regardless of if there is another conflicting one registered after.
