package test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/nbio/st"
)

func TestMockSimple(t *testing.T) {
	defer gock.Off()

	gock.New("http://foo.com").
		Post("/bar").
		MatchType("json").
		Compression("gzip").
		JSON(map[string]string{"foo": "bar"}).
		Reply(201).
		JSON(map[string]string{"bar": "foo"})

	var compressed bytes.Buffer
	w := gzip.NewWriter(&compressed)
	w.Write([]byte(`{"foo":"bar"}`))
	w.Close()
	req, err := http.NewRequest("POST", "http://foo.com/bar", &compressed)
	st.Expect(t, err, nil)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 201)

	resBody, _ := io.ReadAll(res.Body)
	st.Expect(t, string(resBody)[:13], `{"bar":"foo"}`)
}
