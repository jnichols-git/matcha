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

package examples

import (
	"net/http"
	"os"

	"github.com/decentplatforms/matcha/pkg/rctx"
	"github.com/decentplatforms/matcha/pkg/router"
)

type fileServer struct {
	root string
}

func (fs *fileServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := rctx.GetParam(req.Context(), "filepath")
	dat, err := os.ReadFile(fs.root + path)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("File " + path + " does not exist."))
		return
	}
	w.Write(dat)
}

func FileServer(dir string) {
	rt := router.Default()
	rt.Handle(http.MethodGet, "/files/[filepath]+", &fileServer{dir})
	http.ListenAndServe(":3000", rt)
}
