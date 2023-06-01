# Routes

This document details the features of routes.

## Route Expressions

Routes are defined using *expressions* and made up of *parts*. When creating a route, the provided expression will be subdivided along `/`s, and parsed into Parts depending on their syntax. When the route is asked to match a request, the request path will be subdivided in the same way and compared to each Part.

## Route Parts

 Matcha has support for special syntax to customize the behavior of "parts". Currently, it supports:

- Static strings
- Wildcard parameters
- Regex input validation
- Partial routes

### Static Strings

Any part that doesn't match a special part type will be handled as a string literal.

### Wildcards

You can designate a wildcard parameter using a part surrounded by square brackets, `[]`. For example, if you want to have a route that gets data for a device `deviceName`:

```go
r, err := route.New(http.MethodGet, "/devices/[deviceName]/data")
```

### Regex

You can design routes to reject any requests that have improper formatting on a part-by-part basis by using a part surrounded by squiggly brackets, `{}`. You can also combine this with a wildcard to validate route parameters. Building on the following example, if IDs are specifically a string of four lowercase letters or digits in any order:

```go
r, err := route.New(http.MethodGet, "/devices/[deviceName]{[a-z0-9]{4}}/data")
```

Omitting a wildcard parameter will have the same effect; the router just won't store the resulting value.

### Partials

Appending `+` to a route will cause the router to handle requests of a greater or equal length by repeatedly matching against the final part of the route. Using this with a wildcard will set the parameter to the *full additional path*, and using regex will force every part of the additional path to match the regex *individually*.

```go
r, err := route.New(http.MethodGet, "/files/[filename]+")
```

Partials will match against their root with no additional tokens, and if they do, they will not set their parameter.

## Requiring Non-Path Parameters

Matcha provides an interface for matching things that are not paths in package `route/require` as the type `require.Required`. You can define your own with the function definition `func(req *http.Request) bool` and register them onto routes by using the config function or route function `route.Require`.

```go
webRoute, err := route.New(
    http.MethodGet, "/",
    route.Require(require.HostPorts("https://(www.|)cloudretic.com"))
)
apiRoute, err := route.New(
    http.MethodGet(require.HostPorts("https://api.cloudretic.com"))
)
```

We provide some commonly used requirements, such as HostPorts above.

### Hosts

You can use `require.Hosts` or `require.HostPorts` to match hosts using Patterns (regex enclosed in braces). Hosts matches the hostname, while HostPorts matches the scheme, host, and port of the request.
