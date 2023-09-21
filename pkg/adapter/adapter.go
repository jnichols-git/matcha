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

// Package adapter implements an optional interface for use with non-standard requests.
//
// Use of this package is optional when writing adapters for other services, but serves as
// a good representation of what may be needed.
package adapter

import "net/http"

// Adapter is a helper class to define emulation of HTTP requests using arbitrary data.
// It is intended to be fairly permissive, as the use cases for this vary, but any module designed to use
// Matcha for non-HTTP purposes (hosted serverless compute) should use Adapter.
type Adapter[In, Out any] interface {
	// Adapt data In.
	// Must return an http.ResponseWriter, *http.Request pair representing the data In,
	// and a structure Out that is modified by making calls on the ResponseWriter.
	// May return an error if the implementation could potentially fail to Adapt the input,
	// but should not account for errors in handling the resulting request.
	Adapt(In) (http.ResponseWriter, *http.Request, Out, error)
}
