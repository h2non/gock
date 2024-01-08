package gock

import (
	"net/http"

	"github.com/h2non/gock/threadsafe"
)

// var mutex *sync.Mutex = &sync.Mutex{}

var (
	// DefaultTransport stores the default mock transport used by gock.
	DefaultTransport = NewTransport()

	// NativeTransport stores the native net/http default transport
	// in order to restore it when needed.
	NativeTransport = http.DefaultTransport
)

var (
	// ErrCannotMatch store the error returned in case of no matches.
	ErrCannotMatch = threadsafe.ErrCannotMatch
)

// Transport implements http.RoundTripper, which fulfills single http requests issued by
// an http.Client.
//
// gock's Transport encapsulates a given or default http.Transport for further
// delegation, if needed.
type Transport = threadsafe.Transport

// NewTransport creates a new *Transport with no responders.
func NewTransport() *Transport {
	return g.NewTransport(NativeTransport)
}
