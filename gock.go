package gock

import (
	"net/http"
	"sync"

	"github.com/h2non/gock/threadsafe"
)

var g = threadsafe.NewGock()

func init() {
	g.DisableCallback = disable
	g.InterceptCallback = intercept
	g.InterceptingCallback = intercepting
}

// mutex is used interally for locking thread-sensitive functions.
var mutex = &sync.Mutex{}

// ObserverFunc is implemented by users to inspect the outgoing intercepted HTTP traffic
type ObserverFunc = threadsafe.ObserverFunc

// DumpRequest is a default implementation of ObserverFunc that dumps
// the HTTP/1.x wire representation of the http request
var DumpRequest = g.DumpRequest

// New creates and registers a new HTTP mock with
// default settings and returns the Request DSL for HTTP mock
// definition and set up.
func New(uri string) *Request {
	return g.New(uri)
}

// Intercepting returns true if gock is currently able to intercept.
func Intercepting() bool {
	return g.Intercepting()
}

func intercepting() bool {
	mutex.Lock()
	defer mutex.Unlock()
	return http.DefaultTransport == DefaultTransport
}

// Intercept enables HTTP traffic interception via http.DefaultTransport.
// If you are using a custom HTTP transport, you have to use `gock.Transport()`
func Intercept() {
	g.Intercept()
}

func intercept() {
	mutex.Lock()
	http.DefaultTransport = DefaultTransport
	mutex.Unlock()
}

// InterceptClient allows the developer to intercept HTTP traffic using
// a custom http.Client who uses a non default http.Transport/http.RoundTripper implementation.
func InterceptClient(cli *http.Client) {
	g.InterceptClient(cli)
}

// RestoreClient allows the developer to disable and restore the
// original transport in the given http.Client.
func RestoreClient(cli *http.Client) {
	g.RestoreClient(cli)
}

// Disable disables HTTP traffic interception by gock.
func Disable() {
	g.Disable()
}

func disable() {
	mutex.Lock()
	defer mutex.Unlock()
	http.DefaultTransport = NativeTransport
}

// Off disables the default HTTP interceptors and removes
// all the registered mocks, even if they has not been intercepted yet.
func Off() {
	g.Off()
}

// OffAll is like `Off()`, but it also removes the unmatched requests registry.
func OffAll() {
	g.OffAll()
}

// Observe provides a hook to support inspection of the request and matched mock
func Observe(fn ObserverFunc) {
	g.Observe(fn)
}

// EnableNetworking enables real HTTP networking
func EnableNetworking() {
	g.EnableNetworking()
}

// DisableNetworking disables real HTTP networking
func DisableNetworking() {
	g.DisableNetworking()
}

// NetworkingFilter determines if an http.Request should be triggered or not.
func NetworkingFilter(fn FilterRequestFunc) {
	g.NetworkingFilter(fn)
}

// DisableNetworkingFilters disables registered networking filters.
func DisableNetworkingFilters() {
	g.DisableNetworkingFilters()
}

// GetUnmatchedRequests returns all requests that have been received but haven't matched any mock
func GetUnmatchedRequests() []*http.Request {
	return g.GetUnmatchedRequests()
}

// HasUnmatchedRequest returns true if gock has received any requests that didn't match a mock
func HasUnmatchedRequest() bool {
	return g.HasUnmatchedRequest()
}

// CleanUnmatchedRequest cleans the unmatched requests internal registry.
func CleanUnmatchedRequest() {
	g.CleanUnmatchedRequest()
}
