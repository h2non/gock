package threadsafe

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/nbio/st"
)

func TestRegisteredMatchers(t *testing.T) {
	g := NewGock()
	st.Expect(t, len(g.MatchersHeader), 7)
	st.Expect(t, len(g.MatchersBody), 1)
}

func TestNewMatcher(t *testing.T) {
	g := NewGock()
	matcher := g.NewMatcher()
	// Funcs are not comparable, checking slice length as it's better than nothing
	// See https://golang.org/pkg/reflect/#DeepEqual
	st.Expect(t, len(matcher.Matchers), len(g.Matchers))
	st.Expect(t, len(matcher.Get()), len(g.Matchers))
}

func TestNewBasicMatcher(t *testing.T) {
	g := NewGock()
	matcher := g.NewBasicMatcher()
	// Funcs are not comparable, checking slice length as it's better than nothing
	// See https://golang.org/pkg/reflect/#DeepEqual
	st.Expect(t, len(matcher.Matchers), len(g.MatchersHeader))
	st.Expect(t, len(matcher.Get()), len(g.MatchersHeader))
}

func TestNewEmptyMatcher(t *testing.T) {
	g := NewGock()
	matcher := g.NewEmptyMatcher()
	st.Expect(t, len(matcher.Matchers), 0)
	st.Expect(t, len(matcher.Get()), 0)
}

func TestMatcherAdd(t *testing.T) {
	g := NewGock()
	matcher := g.NewMatcher()
	st.Expect(t, len(matcher.Matchers), len(g.Matchers))
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return true, nil
	})
	st.Expect(t, len(matcher.Get()), len(g.Matchers)+1)
}

func TestMatcherSet(t *testing.T) {
	g := NewGock()
	matcher := g.NewMatcher()
	matchers := []MatchFunc{}
	st.Expect(t, len(matcher.Matchers), len(g.Matchers))
	matcher.Set(matchers)
	st.Expect(t, matcher.Matchers, matchers)
	st.Expect(t, len(matcher.Get()), 0)
}

func TestMatcherGet(t *testing.T) {
	g := NewGock()
	matcher := g.NewMatcher()
	matchers := []MatchFunc{}
	matcher.Set(matchers)
	st.Expect(t, matcher.Get(), matchers)
}

func TestMatcherFlush(t *testing.T) {
	g := NewGock()
	matcher := g.NewMatcher()
	st.Expect(t, len(matcher.Matchers), len(g.Matchers))
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return true, nil
	})
	st.Expect(t, len(matcher.Get()), len(g.Matchers)+1)
	matcher.Flush()
	st.Expect(t, len(matcher.Get()), 0)
}

func TestMatcherClone(t *testing.T) {
	g := NewGock()
	matcher := g.DefaultMatcher.Clone()
	st.Expect(t, len(matcher.Get()), len(g.DefaultMatcher.Get()))
}

func TestMatcher(t *testing.T) {
	cases := []struct {
		method  string
		url     string
		matches bool
	}{
		{"GET", "http://foo.com/bar", true},
		{"GET", "http://foo.com/baz", true},
		{"GET", "http://foo.com/foo", false},
		{"POST", "http://foo.com/bar", false},
		{"POST", "http://bar.com/bar", false},
		{"GET", "http://foo.com", false},
	}

	g := NewGock()
	matcher := g.NewMatcher()
	matcher.Flush()
	st.Expect(t, len(matcher.Matchers), 0)

	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return req.Method == "GET", nil
	})
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return req.URL.Host == "foo.com", nil
	})
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return req.URL.Path == "/baz" || req.URL.Path == "/bar", nil
	})

	for _, test := range cases {
		u, _ := url.Parse(test.url)
		req := &http.Request{Method: test.method, URL: u}
		matches, err := matcher.Match(req, nil)
		st.Expect(t, err, nil)
		st.Expect(t, matches, test.matches)
	}
}

func TestMatchMock(t *testing.T) {
	cases := []struct {
		method  string
		url     string
		matches bool
	}{
		{"GET", "http://foo.com/bar", true},
		{"GET", "http://foo.com/baz", true},
		{"GET", "http://foo.com/foo", false},
		{"POST", "http://foo.com/bar", false},
		{"POST", "http://bar.com/bar", false},
		{"GET", "http://foo.com", false},
	}

	g := NewGock()
	matcher := g.DefaultMatcher
	matcher.Flush()
	st.Expect(t, len(matcher.Matchers), 0)

	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return req.Method == "GET", nil
	})
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return req.URL.Host == "foo.com", nil
	})
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return req.URL.Path == "/baz" || req.URL.Path == "/bar", nil
	})

	for _, test := range cases {
		g.Flush()
		mock := g.New(test.url).method(test.method, "").Mock

		u, _ := url.Parse(test.url)
		req := &http.Request{Method: test.method, URL: u}

		match, err := g.MatchMock(req)
		st.Expect(t, err, nil)
		if test.matches {
			st.Expect(t, match, mock)
		} else {
			st.Expect(t, match, nil)
		}
	}

	g.DefaultMatcher.Matchers = g.Matchers
}
