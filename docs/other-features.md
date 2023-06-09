# Other Features

- [Other Features](#other-features)
  - [Cross-Origin Resource Sharing (CORS)](#cross-origin-resource-sharing-cors)
    - [How CORS Works](#how-cors-works)
    - [Setting Up CORS](#setting-up-cors)
    - [Example](#example)
  - [Logging](#logging)
  - [Adapters](#adapters)
    - [Implementing the Adapter Interface](#implementing-the-adapter-interface)
  - [Route Validation](#route-validation)

## Cross-Origin Resource Sharing (CORS)

When requesting resources from a remote server, browsers typically require the server to describe the conditions under which a request may access those resources. This is called Cross-Origin Resource Sharing. Matcha has some tools built in to help you handle CORS requests, if it's required for your application.

### How CORS Works

When a browser sends a request, it first determines if the request is *simple*:

- Uses the `GET` or `POST` method
- Contains no custom headers

If the request is simple, it's sent as normal. If it is not simple, the browser will send a request using the `OPTIONS` method, called a *preflight request*, that asks the server for a set of permitted operations, and if the original request is permitted, sends it as normal. In either case, the browser uses a set of specific headers to determine if access to the resources is permitted.

- `Access-Control-Allow-Origin`: A list of origins that can be used in a request (the URL that the request originates from).
- `Access-Control-Allow-Methods`: A list of methods that can be used in a request.
- `Access-Control-Allow-Headers`: A list of *custom* headers that can be used in a request.
- `Access-Control-Expose-Headers`: A list of *custom* headers that can be accessed by the user agent in the response.
- `Access-Control-Max-Age`: Indicates how long a resource may be cached in seconds.
- `Access-Control-Allow-Credentials`: Indicates if a request may use credentials (cookies, authorization, or TLS).

All of these can be empty, a list, or `*`, which indicates that any value is allowed/exposed. Matcha represents these with the `*AccessControlOptions` struct, used to define how a Router should treat CORS requests.

### Setting Up CORS

There are three ways to set CORS headers on responses.

- `Router` can set the default headers for all routes using the `DefaultCORSHeaders` configuration function.
- `Route` can set the headers for itself only using the `CORSHeaders` configuration function.
- `PreflightCORS` can be used to define an OPTIONS route that returns the given access control headers. *Matcha does not currently automatically generate these routes.*

To manually manipulate CORS headers, `package cors` provides `SetCORSResponseHeaders` that will set the headers based on an `*AccessControlOptions` object. This can be used in the event that the above options don't fit your use case. We'd encourage you to submit an issue on GitHub if your use case isn't immediately supported.

### Example

This router allows all origins, the GET and POST methods, and two custom headers. It will set CORS headers on all responses, and will answer to preflight requests made to `/`.

```go
var aco = &cors.AccessControlOptions{
    AllowOrigin:      []string{"*"},
    AllowMethods:     []string{http.MethodGet, http.MethodPost},
    AllowHeaders:     []string{"X-Some-Header-1", "X-Some-Header-2"},
    MaxAge:           10000,
    AllowCredentials: false,
}
r := Declare(
    Default(),
    DefaultCORS(aco),
    PreflightCORS("/", aco),
    HandleRoute(route.Declare(http.MethodGet, "/"), okHandler("ok")),
    HandleRoute(route.Declare(http.MethodPost, "/"), okHandler("ok")),
    WithNotFound(nfHandler()),
)
```

## Logging

Matcha provides middleware options for logging inbound requests. Each option takes in an `io.Writer`, and writes logs in a specified format for each request. `middleware.LogRequests` and `middleware.LogRequestsIf` will write in the format `[timestamp] [origin] [method] [url]`. Timestamps are in UNIX with nanosecond precision, and the origin will be `-` if it is empty in the request.

## Adapters

Matcha includes the `Adapter` interface to help define functionality that receives HTTP requests through methods other than direct HTTP/S, such as with serverless computing or through a message queue. External adapters aren't required to use this, but it may help.

### Implementing the Adapter Interface

`Adapter` has two type parameters, type `In` and type `Out`, that denote the input and output data being used to emulate an HTTP request. Out should be a pointer type. It has one function, `Adapt`, which takes in a value of type `In` and returns a few values:

- `http.ResponseWriter, *http.Request`: These values should be passed to a `Router` via `ServeHTTP`. The ResponseWriter should write response data to the `Out` value.
- `Out`: This value should be returned in place of normally writing an HTTP response.
- `error`: Adapters may optionally return errors if there could be missing data/incorrect formatting.

## Route Validation

Middleware can be used to validate requests to the router or to specific routes by returning `nil`. Currently, the following are available natively:

- `ExpectQueryParam(name string)` returns 400 Bad Request if a request is missing a query parameter.

Additional validators can be defined using the `middleware.Middleware` type.
