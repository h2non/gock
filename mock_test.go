package gock

import (
	"net/http"
	"testing"

	"github.com/nbio/st"
)

func TestMockSetMatcher(t *testing.T) {
	defer after()

	req := NewRequest()
	res := NewResponse()
	mock := NewMock(req, res)

	matcher := NewMatcher()
	matcher.Flush()
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return true, nil
	})
	mock.SetMatcher(matcher)

	matches, err := mock.Match(&http.Request{})
	st.Expect(t, err, nil)
	st.Expect(t, matches, true)
}

func TestMockAddMatcher(t *testing.T) {
	defer after()

	req := NewRequest()
	res := NewResponse()
	mock := NewMock(req, res)

	matcher := NewMatcher()
	matcher.Flush()
	mock.SetMatcher(matcher)
	mock.AddMatcher(func(req *http.Request, ereq *Request) (bool, error) {
		return true, nil
	})

	matches, err := mock.Match(&http.Request{})
	st.Expect(t, err, nil)
	st.Expect(t, matches, true)
}

func TestMockMatch(t *testing.T) {
	defer after()

	req := NewRequest()
	res := NewResponse()
	mock := NewMock(req, res)

	matcher := NewMatcher()
	matcher.Flush()
	mock.SetMatcher(matcher)
	calls := 0
	mock.AddMatcher(func(req *http.Request, ereq *Request) (bool, error) {
		calls++
		return true, nil
	})
	mock.AddMatcher(func(req *http.Request, ereq *Request) (bool, error) {
		calls++
		return true, nil
	})

	matches, err := mock.Match(&http.Request{})
	st.Expect(t, err, nil)
	st.Expect(t, calls, 2)
	st.Expect(t, matches, true)
}
