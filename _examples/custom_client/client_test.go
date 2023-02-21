package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/nbio/st"
)

func TestClient(t *testing.T) {
	defer gock.Off()

	gock.New("http://foo.com").
		Reply(200).
		BodyString("foo foo")

	req, err := http.NewRequest("GET", "http://foo.com", nil)
	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	res, err := client.Do(req)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 200)
	body, _ := io.ReadAll(res.Body)
	st.Expect(t, string(body), "foo foo")
}
