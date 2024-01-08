package threadsafe

import (
	"bytes"
	"net/http"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/nbio/st"
)

func TestNewRequest(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	req.URL("http://foo.com")
	st.Expect(t, req.URLStruct.Host, "foo.com")
	st.Expect(t, req.URLStruct.Scheme, "http")
	req.MatchHeader("foo", "bar")
	st.Expect(t, req.Header.Get("foo"), "bar")
}

func TestRequestSetURL(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	req.URL("http://foo.com")
	req.SetURL(&url.URL{Host: "bar.com", Path: "/foo"})
	st.Expect(t, req.URLStruct.Host, "bar.com")
	st.Expect(t, req.URLStruct.Path, "/foo")
}

func TestRequestPath(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	req.URL("http://foo.com")
	req.Path("/foo")
	st.Expect(t, req.URLStruct.Scheme, "http")
	st.Expect(t, req.URLStruct.Host, "foo.com")
	st.Expect(t, req.URLStruct.Path, "/foo")
}

func TestRequestBody(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	req.Body(bytes.NewBuffer([]byte("foo bar")))
	st.Expect(t, string(req.BodyBuffer), "foo bar")
}

func TestRequestBodyString(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	req.BodyString("foo bar")
	st.Expect(t, string(req.BodyBuffer), "foo bar")
}

func TestRequestFile(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	absPath, err := filepath.Abs("../version.go")
	st.Expect(t, err, nil)
	req.File(absPath)
	st.Expect(t, string(req.BodyBuffer)[:12], "package gock")
}

func TestRequestJSON(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	req.JSON(map[string]string{"foo": "bar"})
	st.Expect(t, string(req.BodyBuffer)[:13], `{"foo":"bar"}`)
	st.Expect(t, req.Header.Get("Content-Type"), "application/json")
}

func TestRequestXML(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	type xml struct {
		Data string `xml:"data"`
	}
	req.XML(xml{Data: "foo"})
	st.Expect(t, string(req.BodyBuffer), `<xml><data>foo</data></xml>`)
	st.Expect(t, req.Header.Get("Content-Type"), "application/xml")
}

func TestRequestMatchType(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	req.MatchType("json")
	st.Expect(t, req.Header.Get("Content-Type"), "application/json")

	req = g.NewRequest()
	req.MatchType("html")
	st.Expect(t, req.Header.Get("Content-Type"), "text/html")

	req = g.NewRequest()
	req.MatchType("foo/bar")
	st.Expect(t, req.Header.Get("Content-Type"), "foo/bar")
}

func TestRequestBasicAuth(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	req.BasicAuth("bob", "qwerty")
	st.Expect(t, req.Header.Get("Authorization"), "Basic Ym9iOnF3ZXJ0eQ==")
}

func TestRequestMatchHeader(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	req.MatchHeader("foo", "bar")
	req.MatchHeader("bar", "baz")
	req.MatchHeader("UPPERCASE", "bat")
	req.MatchHeader("Mixed-CASE", "foo")

	st.Expect(t, req.Header.Get("foo"), "bar")
	st.Expect(t, req.Header.Get("bar"), "baz")
	st.Expect(t, req.Header.Get("UPPERCASE"), "bat")
	st.Expect(t, req.Header.Get("Mixed-CASE"), "foo")
}

func TestRequestHeaderPresent(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	req.HeaderPresent("foo")
	req.HeaderPresent("bar")
	req.HeaderPresent("UPPERCASE")
	req.HeaderPresent("Mixed-CASE")
	st.Expect(t, req.Header.Get("foo"), ".*")
	st.Expect(t, req.Header.Get("bar"), ".*")
	st.Expect(t, req.Header.Get("UPPERCASE"), ".*")
	st.Expect(t, req.Header.Get("Mixed-CASE"), ".*")
}

func TestRequestMatchParam(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	req.MatchParam("foo", "bar")
	req.MatchParam("bar", "baz")
	st.Expect(t, req.URLStruct.Query().Get("foo"), "bar")
	st.Expect(t, req.URLStruct.Query().Get("bar"), "baz")
}

func TestRequestMatchParams(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	req.MatchParams(map[string]string{"foo": "bar", "bar": "baz"})
	st.Expect(t, req.URLStruct.Query().Get("foo"), "bar")
	st.Expect(t, req.URLStruct.Query().Get("bar"), "baz")
}

func TestRequestPresentParam(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	req.ParamPresent("key")
	st.Expect(t, req.URLStruct.Query().Get("key"), ".*")
}

func TestRequestPathParam(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	req.PathParam("key", "value")
	st.Expect(t, req.PathParams["key"], "value")
}

func TestRequestPersist(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	st.Expect(t, req.Persisted, false)
	req.Persist()
	st.Expect(t, req.Persisted, true)
}

func TestRequestTimes(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	st.Expect(t, req.Counter, 1)
	req.Times(3)
	st.Expect(t, req.Counter, 3)
}

func TestRequestMap(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	st.Expect(t, len(req.Mappers), 0)
	req.Map(func(req *http.Request) *http.Request {
		return req
	})
	st.Expect(t, len(req.Mappers), 1)
}

func TestRequestFilter(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	st.Expect(t, len(req.Filters), 0)
	req.Filter(func(req *http.Request) bool {
		return true
	})
	st.Expect(t, len(req.Filters), 1)
}

func TestRequestEnableNetworking(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	req.Response = &Response{}
	st.Expect(t, req.Response.UseNetwork, false)
	req.EnableNetworking()
	st.Expect(t, req.Response.UseNetwork, true)
}

func TestRequestResponse(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	res := g.NewResponse()
	req.Response = res
	chain := req.Reply(200)
	st.Expect(t, chain, res)
	st.Expect(t, chain.StatusCode, 200)
}

func TestRequestReplyFunc(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	res := g.NewResponse()
	req.Response = res
	chain := req.ReplyFunc(func(r *Response) {
		r.Status(204)
	})
	st.Expect(t, chain, res)
	st.Expect(t, chain.StatusCode, 204)
}

func TestRequestMethods(t *testing.T) {
	g := NewGock()
	req := g.NewRequest()
	req.Get("/foo")
	st.Expect(t, req.Method, "GET")
	st.Expect(t, req.URLStruct.Path, "/foo")

	req = g.NewRequest()
	req.Post("/foo")
	st.Expect(t, req.Method, "POST")
	st.Expect(t, req.URLStruct.Path, "/foo")

	req = g.NewRequest()
	req.Put("/foo")
	st.Expect(t, req.Method, "PUT")
	st.Expect(t, req.URLStruct.Path, "/foo")

	req = g.NewRequest()
	req.Delete("/foo")
	st.Expect(t, req.Method, "DELETE")
	st.Expect(t, req.URLStruct.Path, "/foo")

	req = g.NewRequest()
	req.Patch("/foo")
	st.Expect(t, req.Method, "PATCH")
	st.Expect(t, req.URLStruct.Path, "/foo")

	req = g.NewRequest()
	req.Head("/foo")
	st.Expect(t, req.Method, "HEAD")
	st.Expect(t, req.URLStruct.Path, "/foo")
}

func TestRequestSetMatcher(t *testing.T) {
	g := NewGock()
	defer after(g)

	matcher := g.NewEmptyMatcher()
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return req.URL.Host == "foo.com", nil
	})
	matcher.Add(func(req *http.Request, ereq *Request) (bool, error) {
		return req.Header.Get("foo") == "bar", nil
	})
	ereq := g.NewRequest()
	mock := g.NewMock(ereq, &Response{})
	mock.SetMatcher(matcher)
	ereq.Mock = mock

	headers := make(http.Header)
	headers.Set("foo", "bar")
	req := &http.Request{
		URL:    &url.URL{Host: "foo.com", Path: "/bar"},
		Header: headers,
	}

	match, err := ereq.Mock.Match(req)
	st.Expect(t, err, nil)
	st.Expect(t, match, true)
}

func TestRequestAddMatcher(t *testing.T) {
	g := NewGock()
	defer after(g)

	ereq := g.NewRequest()
	mock := g.NewMock(ereq, &Response{})
	mock.matcher = g.NewMatcher()
	ereq.Mock = mock

	ereq.AddMatcher(func(req *http.Request, ereq *Request) (bool, error) {
		return req.URL.Host == "foo.com", nil
	})
	ereq.AddMatcher(func(req *http.Request, ereq *Request) (bool, error) {
		return req.Header.Get("foo") == "bar", nil
	})

	headers := make(http.Header)
	headers.Set("foo", "bar")
	req := &http.Request{
		URL:    &url.URL{Host: "foo.com", Path: "/bar"},
		Header: headers,
	}

	match, err := ereq.Mock.Match(req)
	st.Expect(t, err, nil)
	st.Expect(t, match, true)
}
