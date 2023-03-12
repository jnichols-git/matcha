# Router Design

This document explains some of the design choices made for Router, why they were made, and what they mean for potential future changes.

## New and Declare

Routes and Routers can be created with either the `New` or `Declare` functions in their respective packages. Their behavior differs slightly in that `New` returns any errors, while `Declare` panics when they occur. The latter behavior is more in line with what you may expect from Go's default router, which handles string paths statically. This was changed due to the impact of increasing the complexity of route registration and handling, which challenges the assumptions of the default router, mainly that *route creation may fail now* if a path is invalid. Idiomatically, you should handle these errors, but historically, Go HTTP routers don't, so we provide both ways.

## Handlers

Some routers use specialized handler functions or types in their APIs. While separating from Go's standard HTTP library can lead to performance improvements, it also leads to stronger coupling with the package being used. We've re-implemented some of the interface members of the HTTP library, allowing for high performance without the need for refactoring code to migrate to `router`. Any handler that works with Go's standard library works for `router` as well.
