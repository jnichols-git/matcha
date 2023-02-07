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

### Creating a Router

You can create a Router with a series of Routes, add Routes after creation, or both, using `WithRoute` and `AddRoute`.

```go
r, err := router.New(
    WithRoute(route.ForceNew("/someEndpoint"), someHandler)
)
r.AddRoute(route.ForceNew("/soeOtherEndpoint"), someOtherHandler)
```

Routes will be handled in the order they are received, and **must match an incoming request URL exactly** in order to call their handler.

You can also add a specific handler that's called in the event that no route matches using `WithNotFound` or `AddNotFound` in the same way (the router will return an empty 404 by default). Doing so will override any previously-set handler for this case.

```go
r, err := router.New(
    WithRoute(...),
    WithNotFound(notFoundHandler)
)
r.AddNotFound(otherNotFoundHandler)
```

If you define custom Middleware, you can attach it to a Router using `Attach`.

```go
router, _ := router.New()
router.Attach(someMiddleware)
```

### Defining Routes

Routes are defined with a string expression delimited by `/`. Creating a route will parse each token into a Part that matches against the token at the same position in incoming requests. There are multiple types of Parts, which are created based on the specific syntax of the token:

- Wildcard: Text enclosed in square brackets `[]`, will match any token at that position and pass the token as a route parameter.
- Regex: Text enclosed in squiggly brackets `{}`, will match any token that is matched **in full** by the contained expression. You should use back-quotes for these routes.
- Static: Any other text.

In the below example, staticRoute will handle requests to `/static`, regexRoute will handle requests to other combination of 5 alphabet letters, and wildcardRoute will handle all other requests (that don't extend beyond that route). The latter two will also store a parameter `word` that contains the value that was matched.

```go
staticRoute, err := route.New("/static")
regexRoute, err := route.New(`/[word]{[a-zA-Z]{5}}`)
wildcardRoute, err := route.New("/[word]")
```

Routes will match GET requests by default; if you want to change that behavior, use `WithMethods`. This will cause the route to no longer match GET requests unless you specify otherwise.

```go
postRoute, err := route.New("/[id]/data", WithMethods(http.MethodPost))
```

If you define custom Middleware, you can attach it to a Route using `Attach`.

```go
route := route.New("/")
route.Attach(someMiddleware)
```

### Middleware

Middleware is defined as  `type Middleware func(*http.Request)*http.Request`. When middleware is attached to a Router or Route, they will be called in-order, and the request will be updated with each one. Middleware can also reject requests by returning `nil`.

Middleware are currently called in the order that they are attached for all implementations of Router and Route.

## Maintainers

Name | Role | Pronouns | GitHub Username | Contact
---|---|---|---|---
Jake Nichols | Creator | they/them | jakenichols2719 | jnichols@cloudretic.com
