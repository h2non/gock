package gock

import (
	"net/http"

	"github.com/h2non/gock/threadsafe"
)

// MatchersHeader exposes an slice of HTTP header specific mock matchers.
func MatchersHeader() []MatchFunc {
	return g.MatchersHeader
}

func SetMatchersHeader(matchers []MatchFunc) {
	g.MatchersHeader = matchers
}

// MatchersBody exposes an slice of HTTP body specific built-in mock matchers.
func MatchersBody() []MatchFunc {
	return g.MatchersBody
}

func SetMatchersBody(matchers []MatchFunc) {
	g.MatchersBody = matchers
}

// Matchers stores all the built-in mock matchers.
func Matchers() []MatchFunc {
	return g.Matchers
}

func SetMatchers(matchers []MatchFunc) {
	g.Matchers = matchers
}

// DefaultMatcher stores the default Matcher instance used to match mocks.
func DefaultMatcher() *MockMatcher {
	return g.DefaultMatcher
}

func SetDefaultMatcher(matcher *MockMatcher) {
	g.DefaultMatcher = matcher
}

// MatchFunc represents the required function
// interface implemented by matchers.
type MatchFunc = threadsafe.MatchFunc

// Matcher represents the required interface implemented by mock matchers.
type Matcher = threadsafe.Matcher

// MockMatcher implements a mock matcher
type MockMatcher = threadsafe.MockMatcher

// NewMatcher creates a new mock matcher
// using the default matcher functions.
func NewMatcher() *MockMatcher {
	return g.NewMatcher()
}

// NewBasicMatcher creates a new matcher with header only mock matchers.
func NewBasicMatcher() *MockMatcher {
	return g.NewBasicMatcher()
}

// NewEmptyMatcher creates a new empty matcher without default matchers.
func NewEmptyMatcher() *MockMatcher {
	return g.NewEmptyMatcher()
}

// MatchMock is a helper function that matches the given http.Request
// in the list of registered mocks, returning it if matches or error if it fails.
func MatchMock(req *http.Request) (Mock, error) {
	return g.MatchMock(req)
}
