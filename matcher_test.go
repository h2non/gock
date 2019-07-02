package gock

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/nbio/st"
)

func TestRegisteredMatchers(t *testing.T) {
	st.Expect(t, len(MatchersHeader), 7)
	st.Expect(t, len(MatchersBody), 1)
}

func TestNewMatcher(t *testing.T) {
	matcher := NewMatcher()
	// Funcs are not comparable, checking slice length as it's better than nothing
	// See https://golang.org/pkg/reflect/#DeepEqual
	st.Expect(t, len(matcher.Matchers), len(Matchers))
	st.Expect(t, len(matcher.Get()), len(Matchers))
}

func TestNewBasicMatcher(t *testing.T) {
	matcher := NewBasicMatcher()
	// Funcs are not comparable, checking slice length as it's better than nothing
	// See https://golang.org/pkg/reflect/#DeepEqual
	st.Expect(t, len(matcher.Matchers), len(MatchersHeader))
	st.Expect(t, len(matcher.Get()), len(MatchersHeader))
}

func TestNewEmptyMatcher(t *testing.T) {
	matcher := NewEmptyMatcher()
	st.Expect(t, len(matcher.Matchers), 0)
	st.Expect(t, len(matcher.Get()), 0)
}

func TestMatcherAdd(t *testing.T) {
	matcher := NewMatcher()
	st.Expect(t, len(matcher.Matchers), len(Matchers))
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return true, nil
	})
	st.Expect(t, len(matcher.Get()), len(Matchers)+1)
}

func TestMatcherSet(t *testing.T) {
	matcher := NewMatcher()
	matchers := []MatchFunc{}
	st.Expect(t, len(matcher.Matchers), len(Matchers))
	matcher.Set(matchers)
	st.Expect(t, matcher.Matchers, matchers)
	st.Expect(t, len(matcher.Get()), 0)
}

func TestMatcherGet(t *testing.T) {
	matcher := NewMatcher()
	matchers := []MatchFunc{}
	matcher.Set(matchers)
	st.Expect(t, matcher.Get(), matchers)
}

func TestMatcherFlush(t *testing.T) {
	matcher := NewMatcher()
	st.Expect(t, len(matcher.Matchers), len(Matchers))
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return true, nil
	})
	st.Expect(t, len(matcher.Get()), len(Matchers)+1)
	matcher.Flush()
	st.Expect(t, len(matcher.Get()), 0)
}

func TestMatcherClone(t *testing.T) {
	matcher := DefaultMatcher.Clone()
	st.Expect(t, len(matcher.Get()), len(DefaultMatcher.Get()))
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

	matcher := NewMatcher()
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

	matcher := DefaultMatcher
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
		Flush()
		mock := New(test.url).method(test.method, "").Mock

		u, _ := url.Parse(test.url)
		req := &http.Request{Method: test.method, URL: u}

		match, err := MatchMock(req)
		st.Expect(t, err, nil)
		if test.matches {
			st.Expect(t, match, mock)
		} else {
			st.Expect(t, match, nil)
		}
	}

	DefaultMatcher.Matchers = Matchers
}
