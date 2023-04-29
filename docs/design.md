# Router Design

This document explains some of the design choices made for Router, why they were made, and what they mean for potential future development.

## New and Declare

Routes and Routers can be created with either the `New` or `Declare` functions in their respective packages. Their behavior differs slightly in that `New` returns any errors, while `Declare` panics when they occur and does not return an error. While [Effective Go](https://go.dev/doc/effective_go#errors) notes that under most circumstances, libraries *should not panic*, the latter method is included for two reasons:

1. Other open source http routers, like the fantastic [httprouter](https://github.com/julienschmidt/httprouter), panic when errors occur during router construction, and we wanted to duplicate that behavior
2. The functional implications of the two things are different; `New` implies the step-by-step creation of routes, while `Declare` implies that the router as defined must exist or the program fails

Generally, it is encouraged that you pick one and stick to it for both routes and routers, to maintain easily predictable behavior.

## Handlers

Some routers use specialized handler functions or types in their APIs. While separating from Go's standard HTTP library can lead to performance improvements, it also leads to stronger coupling with the package being used. We've re-implemented some of the interface members of the HTTP library, allowing for high performance without the need for refactoring code to migrate to Matcha. Any handler that works with Go's standard library works for Matcha as well.
