package threadsafe

import "net/http"

// MatchFunc represents the required function
// interface implemented by matchers.
type MatchFunc func(*http.Request, *Request) (bool, error)

// Matcher represents the required interface implemented by mock matchers.
type Matcher interface {
	// Get returns a slice of registered function matchers.
	Get() []MatchFunc

	// Add adds a new matcher function.
	Add(MatchFunc)

	// Set sets the matchers functions stack.
	Set([]MatchFunc)

	// Flush flushes the current matchers function stack.
	Flush()

	// Match matches the given http.Request with a mock Request.
	Match(*http.Request, *Request) (bool, error)
}

// MockMatcher implements a mock matcher
type MockMatcher struct {
	Matchers []MatchFunc
	g        *Gock
}

// NewMatcher creates a new mock matcher
// using the default matcher functions.
func (g *Gock) NewMatcher() *MockMatcher {
	m := g.NewEmptyMatcher()
	for _, matchFn := range g.Matchers {
		m.Add(matchFn)
	}
	return m
}

// NewBasicMatcher creates a new matcher with header only mock matchers.
func (g *Gock) NewBasicMatcher() *MockMatcher {
	m := g.NewEmptyMatcher()
	for _, matchFn := range g.MatchersHeader {
		m.Add(matchFn)
	}
	return m
}

// NewEmptyMatcher creates a new empty matcher without default matchers.
func (g *Gock) NewEmptyMatcher() *MockMatcher {
	return &MockMatcher{g: g, Matchers: []MatchFunc{}}
}

// Get returns a slice of registered function matchers.
func (m *MockMatcher) Get() []MatchFunc {
	m.g.mutex.Lock()
	defer m.g.mutex.Unlock()
	return m.Matchers
}

// Add adds a new function matcher.
func (m *MockMatcher) Add(fn MatchFunc) {
	m.Matchers = append(m.Matchers, fn)
}

// Set sets a new stack of matchers functions.
func (m *MockMatcher) Set(stack []MatchFunc) {
	m.Matchers = stack
}

// Flush flushes the current matcher
func (m *MockMatcher) Flush() {
	m.Matchers = []MatchFunc{}
}

// Clone returns a separate MockMatcher instance that has a copy of the same MatcherFuncs
func (m *MockMatcher) Clone() *MockMatcher {
	m2 := m.g.NewEmptyMatcher()
	for _, mFn := range m.Get() {
		m2.Add(mFn)
	}
	return m2
}

// Match matches the given http.Request with a mock request
// returning true in case that the request matches, otherwise false.
func (m *MockMatcher) Match(req *http.Request, ereq *Request) (bool, error) {
	for _, matcher := range m.Matchers {
		matches, err := matcher(req, ereq)
		if err != nil {
			return false, err
		}
		if !matches {
			return false, nil
		}
	}
	return true, nil
}

// MatchMock is a helper function that matches the given http.Request
// in the list of registered mocks, returning it if matches or error if it fails.
func (g *Gock) MatchMock(req *http.Request) (Mock, error) {
	for _, mock := range g.GetAll() {
		matches, err := mock.Match(req)
		if err != nil {
			return nil, err
		}
		if matches {
			return mock, nil
		}
	}
	return nil, nil
}
