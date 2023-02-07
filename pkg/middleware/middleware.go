package middleware

import "net/http"

// Middleware runs on any incoming request. Attachment behavior is defined by the structure it's attached to (route vs. router).
//
// Returns an *http.Request; the middleware can set router params or reject a request by returning nil.
type Middleware func(*http.Request) *http.Request
