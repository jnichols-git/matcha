# Cross-Origin Resource Sharing (CORS)

When requesting resources from a remote server, browsers typically require the server to describe the conditions under which a request may access those resources. This is called Cross-Origin Resource Sharing. `router` has some tools built in to help you handle CORS requests, if it's required for your application.

## How CORS Works

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

All of these can be empty, a list, or `*`, which indicates that any value is allowed/exposed. `router` represents these with the `*AccessControlOptions` struct, used to define how a Router should treat CORS requests.

## Setting Up CORS

There are three ways to set CORS headers on responses.

- `Router` can set the default headers for all routes using the `DefaultCORSHeaders` configuration function.
- `Route` can set the headers for itself only using the `CORSHeaders` configuration function.
- `PreflightCORS` can be used to define an OPTIONS route that returns the given access control headers. *`router` does not currently automatically generate these routes.*

To manually manipulate CORS headers, `package cors` provides `SetCORSResponseHeaders` that will set the headers based on an `*AccessControlOptions` object. This can be used in the event that the above options don't fit your use case. We'd encourage you to submit an issue on GitHub if your use case isn't immediately supported.

## Example

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
    WithRoute(route.Declare(http.MethodGet, "/"), okHandler("ok")),
    WithRoute(route.Declare(http.MethodPost, "/"), okHandler("ok")),
    WithNotFound(nfHandler()),
)
```
