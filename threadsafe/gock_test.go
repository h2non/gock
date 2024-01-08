package threadsafe

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/nbio/st"
)

func TestMockSimple(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("http://foo.com").Reply(201).JSON(map[string]string{"foo": "bar"})

	c := &http.Client{}
	g.InterceptClient(c)
	res, err := c.Get("http://foo.com")
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 201)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body)[:13], `{"foo":"bar"}`)
}

func TestMockOff(t *testing.T) {
	g := NewGock()
	g.New("http://foo.com").Reply(201).JSON(map[string]string{"foo": "bar"})
	g.Off()
	c := &http.Client{}
	g.InterceptClient(c)
	_, err := c.Get("http://127.0.0.1:3123")
	st.Reject(t, err, nil)
}

func TestMockBodyStringResponse(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("http://foo.com").Reply(200).BodyString("foo bar")
	c := &http.Client{}
	g.InterceptClient(c)
	res, err := c.Get("http://foo.com")
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 200)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body), "foo bar")
}

func TestMockBodyMatch(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("http://foo.com").BodyString("foo bar").Reply(201).BodyString("foo foo")
	c := &http.Client{}
	g.InterceptClient(c)
	res, err := c.Post("http://foo.com", "text/plain", bytes.NewBuffer([]byte("foo bar")))
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 201)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body), "foo foo")
}

func TestMockBodyCannotMatch(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("http://foo.com").BodyString("foo foo").Reply(201).BodyString("foo foo")
	c := &http.Client{}
	g.InterceptClient(c)
	_, err := c.Post("http://foo.com", "text/plain", bytes.NewBuffer([]byte("foo bar")))
	st.Reject(t, err, nil)
}

func TestMockBodyMatchCompressed(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("http://foo.com").Compression("gzip").BodyString("foo bar").Reply(201).BodyString("foo foo")

	var compressed bytes.Buffer
	w := gzip.NewWriter(&compressed)
	w.Write([]byte("foo bar"))
	w.Close()
	c := &http.Client{}
	g.InterceptClient(c)
	req, err := http.NewRequest("POST", "http://foo.com", &compressed)
	st.Expect(t, err, nil)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "text/plain")
	res, err := c.Do(req)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 201)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body), "foo foo")
}

func TestMockBodyCannotMatchCompressed(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("http://foo.com").Compression("gzip").BodyString("foo bar").Reply(201).BodyString("foo foo")
	c := &http.Client{}
	g.InterceptClient(c)
	_, err := c.Post("http://foo.com", "text/plain", bytes.NewBuffer([]byte("foo bar")))
	st.Reject(t, err, nil)
}

func TestMockBodyMatchJSON(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("http://foo.com").
		Post("/bar").
		JSON(map[string]string{"foo": "bar"}).
		Reply(201).
		JSON(map[string]string{"bar": "foo"})

	c := &http.Client{}
	g.InterceptClient(c)
	res, err := c.Post("http://foo.com/bar", "application/json", bytes.NewBuffer([]byte(`{"foo":"bar"}`)))
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 201)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body)[:13], `{"bar":"foo"}`)
}

func TestMockBodyCannotMatchJSON(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("http://foo.com").
		Post("/bar").
		JSON(map[string]string{"bar": "bar"}).
		Reply(201).
		JSON(map[string]string{"bar": "foo"})

	c := &http.Client{}
	g.InterceptClient(c)
	_, err := c.Post("http://foo.com/bar", "application/json", bytes.NewBuffer([]byte(`{"foo":"bar"}`)))
	st.Reject(t, err, nil)
}

func TestMockBodyMatchCompressedJSON(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("http://foo.com").
		Post("/bar").
		Compression("gzip").
		JSON(map[string]string{"foo": "bar"}).
		Reply(201).
		JSON(map[string]string{"bar": "foo"})

	var compressed bytes.Buffer
	w := gzip.NewWriter(&compressed)
	w.Write([]byte(`{"foo":"bar"}`))
	w.Close()
	c := &http.Client{}
	g.InterceptClient(c)
	req, err := http.NewRequest("POST", "http://foo.com/bar", &compressed)
	st.Expect(t, err, nil)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 201)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body)[:13], `{"bar":"foo"}`)
}

func TestMockBodyCannotMatchCompressedJSON(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("http://foo.com").
		Post("/bar").
		JSON(map[string]string{"bar": "bar"}).
		Reply(201).
		JSON(map[string]string{"bar": "foo"})

	var compressed bytes.Buffer
	w := gzip.NewWriter(&compressed)
	w.Write([]byte(`{"foo":"bar"}`))
	w.Close()
	c := &http.Client{}
	g.InterceptClient(c)
	req, err := http.NewRequest("POST", "http://foo.com/bar", &compressed)
	st.Expect(t, err, nil)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	_, err = c.Do(req)
	st.Reject(t, err, nil)
}

func TestMockMatchHeaders(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("http://foo.com").
		MatchHeader("Content-Type", "(.*)/plain").
		Reply(200).
		BodyString("foo foo")

	c := &http.Client{}
	g.InterceptClient(c)
	res, err := c.Post("http://foo.com", "text/plain", bytes.NewBuffer([]byte("foo bar")))
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 200)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body), "foo foo")
}

func TestMockMap(t *testing.T) {
	g := NewGock()
	defer after(g)

	mock := g.New("http://bar.com")
	mock.Map(func(req *http.Request) *http.Request {
		req.URL.Host = "bar.com"
		return req
	})
	mock.Reply(201).JSON(map[string]string{"foo": "bar"})

	c := &http.Client{}
	g.InterceptClient(c)
	res, err := c.Get("http://foo.com")
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 201)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body)[:13], `{"foo":"bar"}`)
}

func TestMockFilter(t *testing.T) {
	g := NewGock()
	defer after(g)

	mock := g.New("http://foo.com")
	mock.Filter(func(req *http.Request) bool {
		return req.URL.Host == "foo.com"
	})
	mock.Reply(201).JSON(map[string]string{"foo": "bar"})

	c := &http.Client{}
	g.InterceptClient(c)
	res, err := c.Get("http://foo.com")
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 201)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body)[:13], `{"foo":"bar"}`)
}

func TestMockCounterDisabled(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("http://foo.com").Reply(204)
	st.Expect(t, len(g.GetAll()), 1)
	c := &http.Client{}
	g.InterceptClient(c)
	res, err := c.Get("http://foo.com")
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 204)
	st.Expect(t, len(g.GetAll()), 0)
}

func TestMockEnableNetwork(t *testing.T) {
	g := NewGock()
	defer after(g)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world")
	}))
	defer ts.Close()

	g.EnableNetworking()
	defer g.DisableNetworking()

	g.New(ts.URL).Reply(204)
	st.Expect(t, len(g.GetAll()), 1)

	c := &http.Client{}
	g.InterceptClient(c)
	res, err := c.Get(ts.URL)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 204)
	st.Expect(t, len(g.GetAll()), 0)

	res, err = c.Get(ts.URL)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 200)
}

func TestMockEnableNetworkFilter(t *testing.T) {
	g := NewGock()
	defer after(g)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world")
	}))
	defer ts.Close()

	g.EnableNetworking()
	defer g.DisableNetworking()

	g.NetworkingFilter(func(req *http.Request) bool {
		return strings.Contains(req.URL.Host, "127.0.0.1")
	})
	defer g.DisableNetworkingFilters()

	g.New(ts.URL).Reply(0).SetHeader("foo", "bar")
	st.Expect(t, len(g.GetAll()), 1)

	c := &http.Client{}
	g.InterceptClient(c)
	res, err := c.Get(ts.URL)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 200)
	st.Expect(t, res.Header.Get("foo"), "bar")
	st.Expect(t, len(g.GetAll()), 0)
}

func TestMockPersistent(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("http://foo.com").
		Get("/bar").
		Persist().
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	c := &http.Client{}
	g.InterceptClient(c)
	for i := 0; i < 5; i++ {
		res, err := c.Get("http://foo.com/bar")
		st.Expect(t, err, nil)
		st.Expect(t, res.StatusCode, 200)
		body, _ := ioutil.ReadAll(res.Body)
		st.Expect(t, string(body)[:13], `{"foo":"bar"}`)
	}
}

func TestMockPersistTimes(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("http://127.0.0.1:1234").
		Get("/bar").
		Times(4).
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	c := &http.Client{}
	g.InterceptClient(c)
	for i := 0; i < 5; i++ {
		res, err := c.Get("http://127.0.0.1:1234/bar")
		if i == 4 {
			st.Reject(t, err, nil)
			break
		}

		st.Expect(t, err, nil)
		st.Expect(t, res.StatusCode, 200)
		body, _ := ioutil.ReadAll(res.Body)
		st.Expect(t, string(body)[:13], `{"foo":"bar"}`)
	}
}

func TestUnmatched(t *testing.T) {
	g := NewGock()
	defer after(g)

	// clear out any unmatchedRequests from other tests
	g.unmatchedRequests = []*http.Request{}

	g.Intercept()

	c := &http.Client{}
	g.InterceptClient(c)
	_, err := c.Get("http://server.com/unmatched")
	st.Reject(t, err, nil)

	unmatched := g.GetUnmatchedRequests()
	st.Expect(t, len(unmatched), 1)
	st.Expect(t, unmatched[0].URL.Host, "server.com")
	st.Expect(t, unmatched[0].URL.Path, "/unmatched")
	st.Expect(t, g.HasUnmatchedRequest(), true)
}

func TestMultipleMocks(t *testing.T) {
	g := NewGock()
	defer g.Disable()

	g.New("http://server.com").
		Get("/foo").
		Reply(200).
		JSON(map[string]string{"value": "foo"})

	g.New("http://server.com").
		Get("/bar").
		Reply(200).
		JSON(map[string]string{"value": "bar"})

	g.New("http://server.com").
		Get("/baz").
		Reply(200).
		JSON(map[string]string{"value": "baz"})

	tests := []struct {
		path string
	}{
		{"/foo"},
		{"/bar"},
		{"/baz"},
	}

	c := &http.Client{}
	g.InterceptClient(c)
	for _, test := range tests {
		res, err := c.Get("http://server.com" + test.path)
		st.Expect(t, err, nil)
		st.Expect(t, res.StatusCode, 200)
		body, _ := ioutil.ReadAll(res.Body)
		st.Expect(t, string(body)[:15], `{"value":"`+test.path[1:]+`"}`)
	}

	_, err := c.Get("http://server.com/foo")
	st.Reject(t, err, nil)
}

func TestInterceptClient(t *testing.T) {
	g := NewGock()
	defer after(g)

	g.New("http://foo.com").Reply(204)
	st.Expect(t, len(g.GetAll()), 1)

	req, err := http.NewRequest("GET", "http://foo.com", nil)
	client := &http.Client{Transport: &http.Transport{}}
	g.InterceptClient(client)

	res, err := client.Do(req)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 204)
}

func TestRestoreClient(t *testing.T) {
	g := NewGock()
	defer after(g)

	g.New("http://foo.com").Reply(204)
	st.Expect(t, len(g.GetAll()), 1)

	req, err := http.NewRequest("GET", "http://foo.com", nil)
	client := &http.Client{Transport: &http.Transport{}}
	g.InterceptClient(client)
	trans := client.Transport

	res, err := client.Do(req)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 204)

	g.RestoreClient(client)
	st.Reject(t, trans, client.Transport)
}

func TestMockRegExpMatching(t *testing.T) {
	g := NewGock()
	defer after(g)
	g.New("http://foo.com").
		Post("/bar").
		MatchHeader("Authorization", "Bearer (.*)").
		BodyString(`{"foo":".*"}`).
		Reply(200).
		SetHeader("Server", "gock").
		JSON(map[string]string{"foo": "bar"})

	req, _ := http.NewRequest("POST", "http://foo.com/bar", bytes.NewBuffer([]byte(`{"foo":"baz"}`)))
	req.Header.Set("Authorization", "Bearer s3cr3t")

	c := &http.Client{}
	g.InterceptClient(c)
	res, err := c.Do(req)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 200)
	st.Expect(t, res.Header.Get("Server"), "gock")

	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body)[:13], `{"foo":"bar"}`)
}

func TestObserve(t *testing.T) {
	g := NewGock()
	defer after(g)
	var observedRequest *http.Request
	var observedMock Mock
	g.Observe(func(request *http.Request, mock Mock) {
		observedRequest = request
		observedMock = mock
	})
	g.New("http://observe-foo.com").Reply(200)
	req, _ := http.NewRequest("POST", "http://observe-foo.com", nil)

	c := &http.Client{}
	g.InterceptClient(c)
	c.Do(req)

	st.Expect(t, observedRequest.Host, "observe-foo.com")
	st.Expect(t, observedMock.Request().URLStruct.Host, "observe-foo.com")
}

func TestTryCreatingRacesInNew(t *testing.T) {
	g := NewGock()
	defer after(g)
	for i := 0; i < 10; i++ {
		go func() {
			g.New("http://example.com")
		}()
	}
}

func after(g *Gock) {
	g.Flush()
	g.Disable()
}
