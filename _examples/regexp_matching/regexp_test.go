package test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/nbio/st"
)

func TestRegExpMatching(t *testing.T) {
	defer gock.Disable()
	gock.New("http://foo.com").
		Post("/bar").
		MatchHeader("Authorization", "Bearer (.*)").
		BodyString(`{"foo":".*"}`).
		Reply(200).
		SetHeader("Server", "gock").
		JSON(map[string]string{"foo": "bar"})

	req, _ := http.NewRequest("POST", "http://foo.com/bar", bytes.NewBuffer([]byte(`{"foo":"baz"}`)))
	req.Header.Set("Authorization", "Bearer s3cr3t")

	res, err := http.DefaultClient.Do(req)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 200)
	st.Expect(t, res.Header.Get("Server"), "gock")
	body, _ := io.ReadAll(res.Body)
	st.Expect(t, string(body)[:13], `{"foo":"bar"}`)
}
