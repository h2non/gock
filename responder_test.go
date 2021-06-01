package gock

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/nbio/st"
)

func TestResponder(t *testing.T) {
	defer after()
	mres := New("http://foo.com").Reply(200).BodyString("foo")
	req := &http.Request{}

	res, err := Responder(req, mres, nil)
	st.Expect(t, err, nil)
	st.Expect(t, res.Status, "200 OK")
	st.Expect(t, res.StatusCode, 200)

	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body), "foo")
}

func TestResponder_ReadTwice(t *testing.T) {
	defer after()
	mres := New("http://foo.com").Reply(200).BodyString("foo")
	req := &http.Request{}

	res, err := Responder(req, mres, nil)
	st.Expect(t, err, nil)
	st.Expect(t, res.Status, "200 OK")
	st.Expect(t, res.StatusCode, 200)

	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body), "foo")

	body, err = ioutil.ReadAll(res.Body)
	st.Expect(t, err, nil)
	st.Expect(t, body, []byte{})
}

func TestResponderSupportsMultipleHeadersWithSameKey(t *testing.T) {
	defer after()
	mres := New("http://foo").
		Reply(200).
		AddHeader("Set-Cookie", "a=1").
		AddHeader("Set-Cookie", "b=2")
	req := &http.Request{}

	res, err := Responder(req, mres, nil)
	st.Expect(t, err, nil)
	st.Expect(t, res.Header, http.Header{"Set-Cookie": []string{"a=1", "b=2"}})
}

func TestResponderError(t *testing.T) {
	defer after()
	mres := New("http://foo.com").ReplyError(errors.New("error"))
	req := &http.Request{}

	res, err := Responder(req, mres, nil)
	st.Expect(t, err.Error(), "error")
	st.Expect(t, res == nil, true)
}

func TestResponderCancelledContext(t *testing.T) {
	defer after()
	mres := New("http://foo.com").Get("").Reply(200).Delay(20 * time.Millisecond).BodyString("foo")

	// create a context and schedule a call to cancel in 10ms
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://foo.com", nil)

	res, err := Responder(req, mres, nil)

	// verify that we got a context cancellation error and nil response
	st.Expect(t, err, context.Canceled)
	st.Expect(t, res == nil, true)
}

func TestResponderExpiredContext(t *testing.T) {
	defer after()
	mres := New("http://foo.com").Get("").Reply(200).Delay(20 * time.Millisecond).BodyString("foo")

	// create a context that is set to expire in 10ms
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://foo.com", nil)

	res, err := Responder(req, mres, nil)

	// verify that we got a context cancellation error and nil response
	st.Expect(t, err, context.DeadlineExceeded)
	st.Expect(t, res == nil, true)
}

func TestResponderPreExpiredContext(t *testing.T) {
	defer after()
	mres := New("http://foo.com").Get("").Reply(200).BodyString("foo")

	// create a context and wait to ensure it is expired
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Microsecond)
	defer cancel()
	time.Sleep(1 * time.Millisecond)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "http://foo.com", nil)

	res, err := Responder(req, mres, nil)

	// verify that we got a context cancellation error and nil response
	st.Expect(t, err, context.DeadlineExceeded)
	st.Expect(t, res == nil, true)
}
