package cors

import (
	"net/http"

	"github.com/decentplatforms/matcha/pkg/middleware"
)

// CORS middleware
// Assigns the access control options to the related CORS headers to all responses
func CORSMiddleware(aco *AccessControlOptions) middleware.Middleware {
	return func(w http.ResponseWriter, r *http.Request) *http.Request {
		SetCORSResponseHeaders(w, r, aco)
		return r
	}
}
