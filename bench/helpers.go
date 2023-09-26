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

package bench

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/decentplatforms/matcha/pkg/cors"
	"github.com/decentplatforms/matcha/pkg/middleware"
	"github.com/decentplatforms/matcha/pkg/rctx"
	"github.com/decentplatforms/matcha/pkg/route/require"
)

type benchRoute struct {
	method   string
	path     string
	testPath string
	mws      []middleware.Middleware
	rqs      []require.Required
}

func mwCORS() middleware.Middleware {
	return cors.CORSMiddleware(&cors.AccessControlOptions{
		AllowOrigin:  []string{"decentplatforms.com"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{"client_id"},
	})
}

func mwID(w http.ResponseWriter, req *http.Request) *http.Request {
	req.Header.Add("X-Matcha-Request-ID", strconv.FormatInt(rand.Int63(), 10))
	return req
}

func mwIsUserParam(userParam string) middleware.Middleware {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		user := rctx.GetParam(r.Context(), userParam)
		is := r.Header.Get("X-Platform-User-ID")
		if is != user {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("user " + is + " unauthorized"))
			return nil
		}
		return r
	}
}

var api_mws = []middleware.Middleware{mwCORS(), mwID}
var api_mws_auth = []middleware.Middleware{mwIsUserParam("user"), mwCORS(), mwID}
var api_rqs = []require.Required{require.Hosts("{.*}")}
