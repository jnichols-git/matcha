# Routes

- [Route Path Expressions](#route-path-expressions)
  - [Static Strings](#static-strings)
  - [Wildcards](#wildcards)
  - [Regex](#regex)
  - [Partials](#partials)
- [Complex Routes](#complex-routes)
  - [Query Parameters](#query-parameters)
  - [Headers](#headers)
  - [Scheme/Host/Port](#schemehostport)

This document details the features of routes.

## Route Path Expressions

Route paths are defined using *expressions* and made up of *parts*. When creating a route, the provided expression will be subdivided along `/`s, and parsed into Parts depending on their syntax. When the route is asked to match a request, the request path will be subdivided in the same way and compared to each Part.

 Matcha has support for special syntax to customize the behavior of Parts. Currently, it supports:

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

## Complex Routes

You can use `middleware` and `require` to control non-path properties of a request to match against. The most important difference between the two is handling of rejection; `require` will continue checking subsequent routes, while `middleware` will reject the request outright.

### Query Parameters

```go
middleware.ExpectQueryParam(name string, patts ...string)
```

ExpectQueryParam validates that a query parameter is present and matches one of the Patterns provided. If none are provided, any or no value for the parameter is accepted.

### Headers

```go
middleware.ExpectHeader(name string, patts ...string)
```

ExpectHeader validates that a header is present and matches one of the Patterns provided. If none are provided, any value for the header is accepted. Empty header values are invalid by the HTTP spec, so they are not accepted here.

### Scheme/Host/Port

```go
require.Hosts(hostNames ...string)
require.HostPorts(hostNames ...string)
```

Hosts and HostPorts validate that the request was sent to a specific host. The most common use cases for this are routing by scheme (redirect http to https) or subdomain (api vs www). Host names should be provided in the format `[scheme]://[hostname]:[port]`, where:

- `scheme` (optional) is either http or https
- `hostname` is a valid Pattern
- `port` (optional) is a number, number range (inclusive), or comma-delimited list of those two things.

Hosts will only match the `hostname` portion, while HostPorts will match all 3, *even if all 3 are not provided*. `scheme` defaults to http, and `port` defaults to 80 or 443 depending on `scheme`.
