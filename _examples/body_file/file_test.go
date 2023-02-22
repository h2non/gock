package test

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/nbio/st"
)

func TestMockBodyFile(t *testing.T) {
	defer gock.Off()

	gock.New("http://foo.com").
		Post("/bar").
		MatchType("json").
		File("data.json").
		Reply(201).
		File("response.json")

	body := bytes.NewBuffer([]byte(`{"foo":"bar"}`))
	res, err := http.Post("http://foo.com/bar", "application/json", body)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 201)

	resBody, _ := io.ReadAll(res.Body)
	st.Expect(t, string(resBody)[:13], `{"bar":"foo"}`)
}
