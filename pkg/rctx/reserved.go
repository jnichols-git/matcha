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

package rctx

import (
	"context"
)

const PARAM_FULLPATH = string("matcha_fullpath")
const key_reserved_fullpath = paramKey(PARAM_FULLPATH)
const PARAM_MOUNTPROXYTO = string("matcha_mountProxyTo")
const key_reserved_mountProxyTo = paramKey(PARAM_MOUNTPROXYTO)

type reservedParams struct {
	fullpath     string
	mountProxyTo string
}

// get gets a reserved rctx param.
// Returns the value, and true if key is reserved.
func (rps *reservedParams) get(key paramKey) (string, bool) {
	switch key {
	case key_reserved_fullpath:
		return rps.fullpath, true
	case key_reserved_mountProxyTo:
		return rps.mountProxyTo, true
	default:
		return "", false
	}
}

// set sets a reserved rctx param.
// Returns true if the key is reserved.
func (rps *reservedParams) set(parent context.Context, key paramKey, value string) (bool, error) {
	switch key {
	case key_reserved_fullpath:
		if rps.fullpath == "" {
			if parent == nil {
				rps.fullpath = value
			} else if orig := parent.Value(key_reserved_fullpath); orig != nil && orig.(string) != "" {
				rps.fullpath = orig.(string)
			} else {
				rps.fullpath = value
			}
		}
		return true, nil
	case key_reserved_mountProxyTo:
		rps.mountProxyTo = value
		return true, nil
	default:
		return false, nil
	}
}

// reset resets the reserved params.
// Some reserved params have special behavior when set multiple times; this sets back to
// default values so that behavior can be replicated on pooled context.
func (rps *reservedParams) reset() {
	rps.fullpath = ""
	rps.mountProxyTo = ""
}
