// Copyright 2023 Decent Platforms
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

package route

import (
	"github.com/decentplatforms/matcha/pkg/cors"
	"github.com/decentplatforms/matcha/pkg/middleware"
	"github.com/decentplatforms/matcha/pkg/route/require"
)

// RouteConfigFuncs can be applied to a Route at creation.
type ConfigFunc func(Route) error

// Attaches middleware to the route that sets CORS headers on matched requests only.
func CORSHeaders(aco *cors.AccessControlOptions) ConfigFunc {
	return func(r Route) error {
		r.Attach(cors.CORSMiddleware(aco))
		return nil
	}
}

// Attaches middleware to the route.
func WithMiddleware(mws ...middleware.Middleware) ConfigFunc {
	return func(r Route) error {
		r.Attach(mws...)
		return nil
	}
}

func Require(rs ...require.Required) ConfigFunc {
	return func(r Route) error {
		r.Require(rs...)
		return nil
	}
}
