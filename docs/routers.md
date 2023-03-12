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

Routes will be handled in the order they are added.

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
