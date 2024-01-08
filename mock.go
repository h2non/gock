package gock

import (
	"github.com/h2non/gock/threadsafe"
)

// Mock represents the required interface that must
// be implemented by HTTP mock instances.
type Mock = threadsafe.Mock

// Mocker implements a Mock capable interface providing
// a default mock configuration used internally to store mocks.
type Mocker = threadsafe.Mocker

// NewMock creates a new HTTP mock based on the given request and response instances.
// It's mostly used internally.
func NewMock(req *Request, res *Response) *Mocker {
	return g.NewMock(req, res)
}
