# Context in Matcha

In keeping with the "stay fully compatible with the standard library" design goal, Matcha has a custom implementation of `context.Context` in package `rctx`. The primary goal of this package is to provide a more performant *internal* context system. When designing handlers, middleware, requirements, and other structures, you should not need to work with the package; the only necessary public-facing interaction is currently `rctx.GetParam` for fetching URI parameters.

This document serves as a guide for how Matcha handles context internally for developers.

## Attaching Context to a Route

`rctx.PrepareRequestContext` takes in an `http.Request` and a number of parameters to allocate, and returns a request with its context modified in the following ways:

- The old request context becomes the *parent* of the new one.
- The new context has space for `maxParams` params.

## Get/Set Params

`rctx.GetParam` and `rctx.SetParam` are used to manage route parameters (the parts in square brackets), and interact with `*rctx.Context` to store these values. The most common use case for this is just to call `GetParam` with a request context to get a named route parameter.

```go
func HandleReq(w *http.ResponseWriter, req *http.Request) {
    id := rctx.GetParam(req.Context(), "id")
}
```

This works even if the context has been updated in middleware; `GetParam` is type-agnostic, and as long as the original request context is used in the new one, the call will be passed down until the parameter is found or the context chain is exhausted. However, `SetParam` *requires* that the provided context be of type `*rctx.Context` as a safety feature to keep memory use low. As a result, it's recommended that you use `context.WithValue` (or other functions) in middleware instead.
