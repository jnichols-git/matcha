package adapter

import "net/http"

// Adapter is a helper class to define emulation of HTTP requests using arbitrary data.
// It is intended to be fairly permissive, as the use cases for this vary, but any module designed to use
// router for non-HTTP purposes (hosted serverless compute) should use Adapter.
type Adapter[In, Out any] interface {
	// Adapt data In.
	// Must return an http.ResponseWriter, *http.Request pair representing the data In,
	// and a structure Out that is modified by making calls on the ResponseWriter.
	// May return an error if the implementation could potentially fail to Adapt the input,
	// but should not account for errors in handling the resulting request.
	Adapt(In) (http.ResponseWriter, *http.Request, Out, error)
}
