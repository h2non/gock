package threadsafe

import (
	"errors"
	"net/http"
	"sync"
)

var (
	// ErrCannotMatch store the error returned in case of no matches.
	ErrCannotMatch = errors.New("gock: cannot match any request")
)

// Transport implements http.RoundTripper, which fulfills single http requests issued by
// an http.Client.
//
// gock's Transport encapsulates a given or default http.Transport for further
// delegation, if needed.
type Transport struct {
	g *Gock

	// mutex is used to make transport thread-safe of concurrent uses across goroutines.
	mutex sync.Mutex

	// Transport encapsulates the original http.RoundTripper transport interface for delegation.
	Transport http.RoundTripper
}

// NewTransport creates a new *Transport with no responders.
func (g *Gock) NewTransport(transport http.RoundTripper) *Transport {
	return &Transport{g: g, Transport: transport}
}

// transport is used to always return a non-nil transport. This is the same as `(http.Client).transport`, and is what
// would be invoked if gock's transport were not present.
func (m *Transport) transport() http.RoundTripper {
	if m.Transport != nil {
		return m.Transport
	}
	return http.DefaultTransport
}

// RoundTrip receives HTTP requests and routes them to the appropriate responder.  It is required to
// implement the http.RoundTripper interface.  You will not interact with this directly, instead
// the *http.Client you are using will call it for you.
func (m *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Just act as a proxy if not intercepting
	if !m.g.Intercepting() {
		return m.transport().RoundTrip(req)
	}

	m.mutex.Lock()
	defer m.g.Clean()

	var err error
	var res *http.Response

	// Match mock for the incoming http.Request
	mock, err := m.g.MatchMock(req)
	if err != nil {
		m.mutex.Unlock()
		return nil, err
	}

	// Invoke the observer with the intercepted http.Request and matched mock
	if m.g.config.Observer != nil {
		m.g.config.Observer(req, mock)
	}

	// Verify if should use real networking
	networking := shouldUseNetwork(m.g, req, mock)
	if !networking && mock == nil {
		m.mutex.Unlock()
		m.g.trackUnmatchedRequest(req)
		return nil, ErrCannotMatch
	}

	// Ensure me unlock the mutex before building the response
	m.mutex.Unlock()

	// Perform real networking via original transport
	if networking {
		res, err = m.transport().RoundTrip(req)
		// In no mock matched, continue with the response
		if err != nil || mock == nil {
			return res, err
		}
	}

	return Responder(req, mock.Response(), res)
}

// CancelRequest is a no-op function.
func (m *Transport) CancelRequest(req *http.Request) {}

func shouldUseNetwork(g *Gock, req *http.Request, mock Mock) bool {
	if mock != nil && mock.Response().UseNetwork {
		return true
	}
	if !g.config.Networking {
		return false
	}
	if len(g.config.NetworkingFilters) == 0 {
		return true
	}
	for _, filter := range g.config.NetworkingFilters {
		if !filter(req) {
			return false
		}
	}
	return true
}
