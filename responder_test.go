package gock

import (
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

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
