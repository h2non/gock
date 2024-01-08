package gock

import (
	"github.com/h2non/gock/threadsafe"
)

// MapResponseFunc represents the required function interface impletemed by response mappers.
type MapResponseFunc = threadsafe.MapResponseFunc

// FilterResponseFunc represents the required function interface impletemed by response filters.
type FilterResponseFunc = threadsafe.FilterResponseFunc

// Response represents high-level HTTP fields to configure
// and define HTTP responses intercepted by gock.
type Response = threadsafe.Response

// NewResponse creates a new Response.
func NewResponse() *Response {
	return g.NewResponse()
}
