# cloudretic/router

`cloudretic/router` is an actively developed HTTP router for Go, primarily developed for CloudRETIC's API handlers but free to use by anyone under the Apache 2.0 license.

## Installation

`go get github.com/cloudretic/router`

## Usage

### Components

`router` has three main components:

- Routers implement `http.Handler` and handle incoming requests using a combination of Routes and Middleware.
- Routes are derived from a string expression and match against `*http.Request`s.
- Middleware acts on `*http.Requests`, usually to either transform or reject them.

### New vs. Declare

Routes and Routers can both be created via the package function `New` or `Declare`. `New` returns the object and an error (if one occurs), while `Declare` only returns the object, and will panic if creation fails. It is recommended that you use `New` if you need to create or configure a router at runtime, and `Declare` if you're creating a static router when the program starts.

```go
infoRoute, err := route.New("/info")
if err != nil {
    ...
}
fileRoute, err := route.New("/file/[filePath]+")
if err != nil {
    ...
}
rt, err := router.New(
    WithRoute(infoRoute, infoHandler),
    WithRoute(fileRoute, fileHandler),
)
if err != nil {
    ...
}
```

```go
// Panics if there's an error creating the router
rt := router.Declare(
    WithRoute(route.Declare("/info"), infoHandler),
    WithRoute(route.Declare("/file/[filePath]+"), fileHandler)
)
```

### Middleware

Middleware is defined as  `type Middleware func(*http.Request)*http.Request`. When middleware is attached to a Router or Route, they will be called in-order, and the request will be updated with each one. Middleware can also reject requests by returning `nil`.

Middleware are currently called in the order that they are attached for all implementations of Router and Route.

## Maintainers

Name | Role | Pronouns | GitHub Username | Contact
---|---|---|---|---
Jake Nichols | Creator | they/them | jakenichols2719 | jnichols@cloudretic.com
