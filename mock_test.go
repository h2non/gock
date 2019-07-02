package gock

import (
	"net/http"
	"testing"

	"github.com/nbio/st"
)

func TestNewMock(t *testing.T) {
	defer after()

	req := NewRequest()
	res := NewResponse()
	mock := NewMock(req, res)
	st.Expect(t, mock.disabled, false)
	st.Expect(t, len(mock.matcher.Get()), len(DefaultMatcher.Get()))

	st.Expect(t, mock.Request(), req)
	st.Expect(t, mock.Request().Mock, mock)
	st.Expect(t, mock.Response(), res)
	st.Expect(t, mock.Response().Mock, mock)
}

func TestMockDisable(t *testing.T) {
	defer after()

	req := NewRequest()
	res := NewResponse()
	mock := NewMock(req, res)

	st.Expect(t, mock.disabled, false)
	mock.Disable()
	st.Expect(t, mock.disabled, true)

	matches, err := mock.Match(&http.Request{})
	st.Expect(t, err, nil)
	st.Expect(t, matches, false)
}

func TestMockDone(t *testing.T) {
	defer after()

	req := NewRequest()
	res := NewResponse()

	mock := NewMock(req, res)
	st.Expect(t, mock.disabled, false)
	st.Expect(t, mock.Done(), false)

	mock = NewMock(req, res)
	st.Expect(t, mock.disabled, false)
	mock.Disable()
	st.Expect(t, mock.Done(), true)

	mock = NewMock(req, res)
	st.Expect(t, mock.disabled, false)
	mock.request.Counter = 0
	st.Expect(t, mock.Done(), true)

	mock = NewMock(req, res)
	st.Expect(t, mock.disabled, false)
	mock.request.Persisted = true
	st.Expect(t, mock.Done(), false)
}

func TestMockSetMatcher(t *testing.T) {
	defer after()

	req := NewRequest()
	res := NewResponse()
	mock := NewMock(req, res)

	st.Expect(t, len(mock.matcher.Get()), len(DefaultMatcher.Get()))
	matcher := NewMatcher()
	matcher.Flush()
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return true, nil
	})
	mock.SetMatcher(matcher)
	st.Expect(t, len(mock.matcher.Get()), 1)
	st.Expect(t, mock.disabled, false)

	matches, err := mock.Match(&http.Request{})
	st.Expect(t, err, nil)
	st.Expect(t, matches, true)
}

func TestMockAddMatcher(t *testing.T) {
	defer after()

	req := NewRequest()
	res := NewResponse()
	mock := NewMock(req, res)

	st.Expect(t, len(mock.matcher.Get()), len(DefaultMatcher.Get()))
	matcher := NewMatcher()
	matcher.Flush()
	mock.SetMatcher(matcher)
	mock.AddMatcher(func(req *http.Request, ereq *Request) (bool, error) {
		return true, nil
	})
	st.Expect(t, mock.disabled, false)
	st.Expect(t, mock.matcher, matcher)

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
	st.Expect(t, mock.disabled, false)
	st.Expect(t, mock.matcher, matcher)

	matches, err := mock.Match(&http.Request{})
	st.Expect(t, err, nil)
	st.Expect(t, calls, 2)
	st.Expect(t, matches, true)
}
