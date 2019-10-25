package gock

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/nbio/st"
)

func TestMatchMethod(t *testing.T) {
	cases := []struct {
		value   string
		method  string
		matches bool
	}{
		{"GET", "GET", true},
		{"POST", "POST", true},
		{"", "POST", true},
		{"POST", "GET", false},
		{"PUT", "GET", false},
	}

	for _, test := range cases {
		req := &http.Request{Method: test.method}
		ereq := &Request{Method: test.value}
		matches, err := MatchMethod(req, ereq)
		st.Expect(t, err, nil)
		st.Expect(t, matches, test.matches)
	}
}

func TestMatchScheme(t *testing.T) {
	cases := []struct {
		value   string
		scheme  string
		matches bool
	}{
		{"http", "http", true},
		{"https", "https", true},
		{"http", "https", false},
		{"", "https", true},
		{"https", "", true},
	}

	for _, test := range cases {
		req := &http.Request{URL: &url.URL{Scheme: test.scheme}}
		ereq := &Request{URLStruct: &url.URL{Scheme: test.value}}
		matches, err := MatchScheme(req, ereq)
		st.Expect(t, err, nil)
		st.Expect(t, matches, test.matches)
	}
}

func TestMatchHost(t *testing.T) {
	cases := []struct {
		value            string
		url              string
		matches          bool
		matchesNonRegexp bool
	}{
		{"foo.com", "foo.com", true, true},
		{"FOO.com", "foo.com", true, true},
		{"foo.net", "foo.com", false, false},
		{"foo.bar.net", "foo-bar.net", true, false},
		{"foo", "foo.com", true, false},
		{"(.*).com", "foo.com", true, false},
		{"127.0.0.1", "127.0.0.1", true, true},
		{"127.0.0.2", "127.0.0.1", false, false},
		{"127.0.0.*", "127.0.0.1", true, false},
		{"127.0.0.[0-9]", "127.0.0.7", true, false},
	}

	for _, test := range cases {
		req := &http.Request{URL: &url.URL{Host: test.url}}
		ereq := &Request{URLStruct: &url.URL{Host: test.value}}
		matches, err := MatchHost(req, ereq)
		st.Expect(t, err, nil)
		st.Expect(t, matches, test.matches)
		ereq.WithOptions(Options{DisableRegexpHost: true})
		matches, err = MatchHost(req, ereq)
		st.Expect(t, err, nil)
		st.Expect(t, matches, test.matchesNonRegexp)
	}
}

func TestMatchPath(t *testing.T) {
	cases := []struct {
		value   string
		path    string
		matches bool
	}{
		{"/foo", "/foo", true},
		{"/foo", "/foo/bar", true},
		{"bar", "/foo/bar", true},
		{"foo", "/foo/bar", true},
		{"bar$", "/foo/bar", true},
		{"/foo/*", "/foo/bar", true},
		{"/foo/[a-z]+", "/foo/bar", true},
		{"/foo/baz", "/foo/bar", false},
		{"/foo/baz", "/foo/bar", false},
		{"/foo/bar%3F+%C3%A9", "/foo/bar%3F+%C3%A9", true},
	}

	for _, test := range cases {
		u, _ := url.Parse("http://foo.com" + test.path)
		mu, _ := url.Parse("http://foo.com" + test.value)
		req := &http.Request{URL: u}
		ereq := &Request{URLStruct: mu}
		matches, err := MatchPath(req, ereq)
		st.Expect(t, err, nil)
		st.Expect(t, matches, test.matches)
	}
}

func TestMatchHeaders(t *testing.T) {
	cases := []struct {
		values  http.Header
		headers http.Header
		matches bool
	}{
		{http.Header{"foo": []string{"bar"}}, http.Header{"foo": []string{"bar"}}, true},
		{http.Header{"foo": []string{"bar"}}, http.Header{"foo": []string{"barbar"}}, true},
		{http.Header{"bar": []string{"bar"}}, http.Header{"foo": []string{"bar"}}, false},
		{http.Header{"foofoo": []string{"bar"}}, http.Header{"foo": []string{"bar"}}, false},
		{http.Header{"foo": []string{"bar(.*)"}}, http.Header{"foo": []string{"barbar"}}, true},
		{http.Header{"foo": []string{"b(.*)"}}, http.Header{"foo": []string{"barbar"}}, true},
		{http.Header{"foo": []string{"^bar$"}}, http.Header{"foo": []string{"bar"}}, true},
		{http.Header{"foo": []string{"^bar$"}}, http.Header{"foo": []string{"barbar"}}, false},
		{http.Header{"UPPERCASE": []string{"bar"}}, http.Header{"UPPERCASE": []string{"bar"}}, true},
		{http.Header{"Mixed-CASE": []string{"bar"}}, http.Header{"Mixed-CASE": []string{"bar"}}, true},
		{http.Header{"User-Agent": []string{"Agent (version1.0)"}}, http.Header{"User-Agent": []string{"Agent (version1.0)"}}, true},
		{http.Header{"Content-Type": []string{"(.*)/plain"}}, http.Header{"Content-Type": []string{"text/plain"}}, true},
	}

	for _, test := range cases {
		req := &http.Request{Header: test.headers}
		ereq := &Request{Header: test.values}
		matches, err := MatchHeaders(req, ereq)
		st.Expect(t, err, nil)
		st.Expect(t, matches, test.matches)
	}
}

func TestMatchQueryParams(t *testing.T) {
	cases := []struct {
		value   string
		path    string
		matches bool
	}{
		{"foo=bar", "foo=bar", true},
		{"foo=bar", "foo=foo&foo=bar", true},
		{"foo=b*", "foo=bar", true},
		{"foo=.*", "foo=bar", true},
		{"foo=f[o]{2}", "foo=foo", true},
		{"foo=bar&bar=foo", "foo=bar&foo=foo&bar=foo", true},
		{"foo=", "foo=bar", true},
		{"foo=foo", "foo=bar", false},
		{"bar=bar", "foo=bar bar", false},
	}

	for _, test := range cases {
		u, _ := url.Parse("http://foo.com/?" + test.path)
		mu, _ := url.Parse("http://foo.com/?" + test.value)
		req := &http.Request{URL: u}
		ereq := &Request{URLStruct: mu}
		matches, err := MatchQueryParams(req, ereq)
		st.Expect(t, err, nil)
		st.Expect(t, matches, test.matches)
	}
}

func TestMatchPathParams(t *testing.T) {
	cases := []struct {
		key     string
		value   string
		path    string
		matches bool
	}{
		{"foo", "bar", "/foo/bar", true},
		{"foo", "bar", "/foo/test/bar", false},
		{"foo", "bar", "/test/foo/bar/ack", true},
		{"foo", "bar", "/foo", false},
	}

	for i, test := range cases {
		u, _ := url.Parse("http://foo.com" + test.path)
		mu, _ := url.Parse("http://foo.com" + test.path)
		req := &http.Request{URL: u}
		ereq := &Request{
			URLStruct:  mu,
			PathParams: map[string]string{test.key: test.value},
		}
		matches, err := MatchPathParams(req, ereq)
		st.Expect(t, err, nil, i)
		st.Expect(t, matches, test.matches, i)
	}
}

func TestMatchBody(t *testing.T) {
	cases := []struct {
		value   string
		body    string
		matches bool
	}{
		{"foo bar", "foo bar\n", true},
		{"foo", "foo bar\n", true},
		{"f[o]+", "foo\n", true},
		{`"foo"`, `{"foo":"bar"}\n`, true},
		{`{"foo":"bar"}`, `{"foo":"bar"}\n`, true},
		{`{"foo":"foo"}`, `{"foo":"bar"}\n`, false},

		{`{"foo":"bar","bar":"foo"}`, `{"bar":"foo","foo":"bar"}`, true},
		{`{"bar":"foo","foo":{"two":"three","three":"two"}}`, `{"foo":{"three":"two","two":"three"},"bar":"foo"}`, true},
	}

	for _, test := range cases {
		req := &http.Request{Body: createReadCloser([]byte(test.body))}
		ereq := &Request{BodyBuffer: []byte(test.value)}
		matches, err := MatchBody(req, ereq)
		st.Expect(t, err, nil)
		st.Expect(t, matches, test.matches)
	}
}
