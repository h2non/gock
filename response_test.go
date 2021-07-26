package gock

import (
	"bytes"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/nbio/st"
)

func TestNewResponse(t *testing.T) {
	res := NewResponse()

	res.Status(200)
	st.Expect(t, res.StatusCode, 200)

	res.SetHeader("foo", "bar")
	st.Expect(t, res.Header.Get("foo"), "bar")

	res.Delay(1000 * time.Millisecond)
	st.Expect(t, res.ResponseDelay, 1000*time.Millisecond)

	res.EnableNetworking()
	st.Expect(t, res.UseNetwork, true)
}

func TestResponseStatus(t *testing.T) {
	res := NewResponse()
	st.Expect(t, res.StatusCode, 0)
	res.Status(200)
	st.Expect(t, res.StatusCode, 200)
}

func TestResponseType(t *testing.T) {
	res := NewResponse()
	res.Type("json")
	st.Expect(t, res.Header.Get("Content-Type"), "application/json")

	res = NewResponse()
	res.Type("xml")
	st.Expect(t, res.Header.Get("Content-Type"), "application/xml")

	res = NewResponse()
	res.Type("foo/bar")
	st.Expect(t, res.Header.Get("Content-Type"), "foo/bar")
}

func TestResponseSetHeader(t *testing.T) {
	res := NewResponse()
	res.SetHeader("foo", "bar")
	res.SetHeader("bar", "baz")
	st.Expect(t, res.Header.Get("foo"), "bar")
	st.Expect(t, res.Header.Get("bar"), "baz")
}

func TestResponseAddHeader(t *testing.T) {
	res := NewResponse()
	res.AddHeader("foo", "bar")
	res.AddHeader("foo", "baz")
	st.Expect(t, res.Header.Get("foo"), "bar")
	st.Expect(t, res.Header["Foo"][1], "baz")
}

func TestResponseSetHeaders(t *testing.T) {
	res := NewResponse()
	res.SetHeaders(map[string]string{"foo": "bar", "bar": "baz"})
	st.Expect(t, res.Header.Get("foo"), "bar")
	st.Expect(t, res.Header.Get("bar"), "baz")
}

func TestResponseBody(t *testing.T) {
	res := NewResponse()
	res.Body(bytes.NewBuffer([]byte("foo bar")))
	st.Expect(t, string(res.BodyBuffer), "foo bar")
}

func TestResponseBodyString(t *testing.T) {
	res := NewResponse()
	res.BodyString("foo bar")
	st.Expect(t, string(res.BodyBuffer), "foo bar")
}

func TestResponseFile(t *testing.T) {
	res := NewResponse()
	res.File("version.go")
	st.Expect(t, string(res.BodyBuffer)[:12], "package gock")
}

func TestResponseJSON(t *testing.T) {
	res := NewResponse()
	res.JSON(map[string]string{"foo": "bar"})
	st.Expect(t, string(res.BodyBuffer)[:13], `{"foo":"bar"}`)
	st.Expect(t, res.Header.Get("Content-Type"), "application/json")
}

func TestResponseXML(t *testing.T) {
	res := NewResponse()
	type xml struct {
		Data string `xml:"data"`
	}
	res.XML(xml{Data: "foo"})
	st.Expect(t, string(res.BodyBuffer), `<xml><data>foo</data></xml>`)
	st.Expect(t, res.Header.Get("Content-Type"), "application/xml")
}

func TestResponseMap(t *testing.T) {
	res := NewResponse()
	st.Expect(t, len(res.Mappers), 0)
	res.Map(func(res *http.Response) *http.Response {
		return res
	})
	st.Expect(t, len(res.Mappers), 1)
}

func TestResponseFilter(t *testing.T) {
	res := NewResponse()
	st.Expect(t, len(res.Filters), 0)
	res.Filter(func(res *http.Response) bool {
		return true
	})
	st.Expect(t, len(res.Filters), 1)
}

func TestResponseSetError(t *testing.T) {
	res := NewResponse()
	st.Expect(t, res.Error, nil)
	res.SetError(errors.New("foo error"))
	st.Expect(t, res.Error.Error(), "foo error")
}

func TestResponseDelay(t *testing.T) {
	res := NewResponse()
	st.Expect(t, res.ResponseDelay, 0*time.Microsecond)
	res.Delay(100 * time.Millisecond)
	st.Expect(t, res.ResponseDelay, 100*time.Millisecond)
}

func TestResponseEnableNetworking(t *testing.T) {
	res := NewResponse()
	st.Expect(t, res.UseNetwork, false)
	res.EnableNetworking()
	st.Expect(t, res.UseNetwork, true)
}

func TestResponseDone(t *testing.T) {
	res := NewResponse()
	res.Mock = &Mocker{request: &Request{Counter: 1}, disabler: new(disabler)}
	st.Expect(t, res.Done(), false)
	res.Mock.Disable()
	st.Expect(t, res.Done(), true)
}
