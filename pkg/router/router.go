// Copyright 2023 Matcha Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package router defines the API for creating and interacting with a router, and its internal behavior.
//
// See [https://github.com/decentplatforms/matcha/blob/main/docs/routers.md].
package router

import (
	"net/http"

	"github.com/decentplatforms/matcha/pkg/middleware"
	"github.com/decentplatforms/matcha/pkg/route"
)

type Router interface {
	// Attach middleware to the router.
	//
	// Router implementations must run all attached middleware on all incoming requests.
	Attach(mw ...middleware.Middleware)
	// Add a route to the router.
	//
	// AddRoute was deprecated in v1.2.0. Use HandleRoute instead.
	AddRoute(r route.Route, h http.Handler)
	// Handle a method and path.
	// This constructs a basic Route internally. Returns an error if routing path rules are
	// violated; see routes.md.
	Handle(method, path string, h http.Handler) error
	// Handle a method and path.
	// This constructs a basic Route internally. Returns an error if routing path rules are
	// violated; see routes.md.
	HandleFunc(method, path string, h http.HandlerFunc) error
	// Handle a more complex path.
	// If you're only using method+path, use Handle instead.
	HandleRoute(r route.Route, h http.Handler)
	// Handle a more complex path.
	// If you're only using method+path, use Handle instead.
	HandleRouteFunc(r route.Route, h http.HandlerFunc)
	// Mount a handler at a path.
	// Forwards all requests at path to the provided handler, optionally limited to a set
	// of methods passed in the variadic methods parameter. Use this if you want to
	// use your existing handler at a specific URI.
	Mount(path string, h http.Handler, methods ...string) error
	// Add a handler for any request that is not matched.
	//
	// Router implementations should define default behavior, and must allow user assignment of behavior.
	AddNotFound(h http.Handler)
	// Implements http.Handler
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// Create a new Router.
// Returns an error if creation fails.
func New(with Router, cfs ...ConfigFunc) (Router, error) {
	for _, cf := range cfs {
		err := cf(with)
		if err != nil {
			return nil, err
		}
	}
	return with, nil
}

// Declare a new Router.
// Panics if creation fails.
func Declare(with Router, cfs ...ConfigFunc) Router {
	for _, cf := range cfs {
		err := cf(with)
		if err != nil {
			panic(err)
		}
	}
	return with
}
