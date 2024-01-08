package threadsafe

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"sync"
)

type Gock struct {
	// mutex is used internally for locking thread-sensitive functions.
	mutex sync.Mutex
	// config global singleton store.
	config struct {
		Networking        bool
		NetworkingFilters []FilterRequestFunc
		Observer          ObserverFunc
	}
	// DumpRequest is a default implementation of ObserverFunc that dumps
	// the HTTP/1.x wire representation of the http request
	DumpRequest ObserverFunc
	// track unmatched requests so they can be tested for
	unmatchedRequests []*http.Request

	// storeMutex is used internally for store synchronization.
	storeMutex sync.RWMutex

	// mocks is internally used to store registered mocks.
	mocks []Mock

	// DefaultMatcher stores the default Matcher instance used to match mocks.
	DefaultMatcher *MockMatcher

	// MatchersHeader exposes a slice of HTTP header specific mock matchers.
	MatchersHeader []MatchFunc
	// MatchersBody exposes a slice of HTTP body specific built-in mock matchers.
	MatchersBody []MatchFunc
	// Matchers stores all the built-in mock matchers.
	Matchers []MatchFunc

	// BodyTypes stores the supported MIME body types for matching.
	// Currently only text-based types.
	BodyTypes []string

	// BodyTypeAliases stores a generic MIME type by alias.
	BodyTypeAliases map[string]string

	// CompressionSchemes stores the supported Content-Encoding types for decompression.
	CompressionSchemes []string

	intercepting bool

	DisableCallback      func()
	InterceptCallback    func()
	InterceptingCallback func() bool
}

func NewGock() *Gock {
	g := &Gock{
		DumpRequest: defaultDumpRequest,

		BodyTypes: []string{
			"text/html",
			"text/plain",
			"application/json",
			"application/xml",
			"multipart/form-data",
			"application/x-www-form-urlencoded",
		},

		BodyTypeAliases: map[string]string{
			"html": "text/html",
			"text": "text/plain",
			"json": "application/json",
			"xml":  "application/xml",
			"form": "multipart/form-data",
			"url":  "application/x-www-form-urlencoded",
		},

		// CompressionSchemes stores the supported Content-Encoding types for decompression.
		CompressionSchemes: []string{
			"gzip",
		},
	}
	g.MatchersHeader = []MatchFunc{
		g.MatchMethod,
		g.MatchScheme,
		g.MatchHost,
		g.MatchPath,
		g.MatchHeaders,
		g.MatchQueryParams,
		g.MatchPathParams,
	}
	g.MatchersBody = []MatchFunc{
		g.MatchBody,
	}
	g.Matchers = append(g.MatchersHeader, g.MatchersBody...)

	// DefaultMatcher stores the default Matcher instance used to match mocks.
	g.DefaultMatcher = g.NewMatcher()
	return g
}

// ObserverFunc is implemented by users to inspect the outgoing intercepted HTTP traffic
type ObserverFunc func(*http.Request, Mock)

func defaultDumpRequest(request *http.Request, mock Mock) {
	bytes, _ := httputil.DumpRequestOut(request, true)
	fmt.Println(string(bytes))
	fmt.Printf("\nMatches: %v\n---\n", mock != nil)
}

// New creates and registers a new HTTP mock with
// default settings and returns the Request DSL for HTTP mock
// definition and set up.
func (g *Gock) New(uri string) *Request {
	g.Intercept()

	res := g.NewResponse()
	req := g.NewRequest()
	req.URLStruct, res.Error = url.Parse(normalizeURI(uri))

	// Create the new mock expectation
	exp := g.NewMock(req, res)
	g.Register(exp)

	return req
}

// Intercepting returns true if gock is currently able to intercept.
func (g *Gock) Intercepting() bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	callbackResponse := true
	if g.InterceptingCallback != nil {
		callbackResponse = g.InterceptingCallback()
	}

	return g.intercepting && callbackResponse
}

// Intercept enables HTTP traffic interception via http.DefaultTransport.
// If you are using a custom HTTP transport, you have to use `gock.Transport()`
func (g *Gock) Intercept() {
	if !g.Intercepting() {
		g.mutex.Lock()
		g.intercepting = true

		if g.InterceptCallback != nil {
			g.InterceptCallback()
		}

		g.mutex.Unlock()
	}
}

// InterceptClient allows the developer to intercept HTTP traffic using
// a custom http.Client who uses a non default http.Transport/http.RoundTripper implementation.
func (g *Gock) InterceptClient(cli *http.Client) {
	_, ok := cli.Transport.(*Transport)
	if ok {
		return // if transport already intercepted, just ignore it
	}
	cli.Transport = g.NewTransport(cli.Transport)
}

// RestoreClient allows the developer to disable and restore the
// original transport in the given http.Client.
func (g *Gock) RestoreClient(cli *http.Client) {
	trans, ok := cli.Transport.(*Transport)
	if !ok {
		return
	}
	cli.Transport = trans.Transport
}

// Disable disables HTTP traffic interception by gock.
func (g *Gock) Disable() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.intercepting = false

	if g.DisableCallback != nil {
		g.DisableCallback()
	}
}

// Off disables the default HTTP interceptors and removes
// all the registered mocks, even if they has not been intercepted yet.
func (g *Gock) Off() {
	g.Flush()
	g.Disable()
}

// OffAll is like `Off()`, but it also removes the unmatched requests registry.
func (g *Gock) OffAll() {
	g.Flush()
	g.Disable()
	g.CleanUnmatchedRequest()
}

// Observe provides a hook to support inspection of the request and matched mock
func (g *Gock) Observe(fn ObserverFunc) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.config.Observer = fn
}

// EnableNetworking enables real HTTP networking
func (g *Gock) EnableNetworking() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.config.Networking = true
}

// DisableNetworking disables real HTTP networking
func (g *Gock) DisableNetworking() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.config.Networking = false
}

// NetworkingFilter determines if an http.Request should be triggered or not.
func (g *Gock) NetworkingFilter(fn FilterRequestFunc) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.config.NetworkingFilters = append(g.config.NetworkingFilters, fn)
}

// DisableNetworkingFilters disables registered networking filters.
func (g *Gock) DisableNetworkingFilters() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.config.NetworkingFilters = []FilterRequestFunc{}
}

// GetUnmatchedRequests returns all requests that have been received but haven't matched any mock
func (g *Gock) GetUnmatchedRequests() []*http.Request {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	return g.unmatchedRequests
}

// HasUnmatchedRequest returns true if gock has received any requests that didn't match a mock
func (g *Gock) HasUnmatchedRequest() bool {
	return len(g.GetUnmatchedRequests()) > 0
}

// CleanUnmatchedRequest cleans the unmatched requests internal registry.
func (g *Gock) CleanUnmatchedRequest() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.unmatchedRequests = []*http.Request{}
}

func (g *Gock) trackUnmatchedRequest(req *http.Request) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.unmatchedRequests = append(g.unmatchedRequests, req)
}

func normalizeURI(uri string) string {
	if ok, _ := regexp.MatchString("^http[s]?", uri); !ok {
		return "http://" + uri
	}
	return uri
}
