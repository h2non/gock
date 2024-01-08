package gock

import (
	"github.com/h2non/gock/threadsafe"
)

// MapRequestFunc represents the required function interface for request mappers.
type MapRequestFunc = threadsafe.MapRequestFunc

// FilterRequestFunc represents the required function interface for request filters.
type FilterRequestFunc = threadsafe.FilterRequestFunc

// Request represents the high-level HTTP request used to store
// request fields used to match intercepted requests.
type Request = threadsafe.Request

// NewRequest creates a new Request instance.
func NewRequest() *Request {
	return g.NewRequest()
}
