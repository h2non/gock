package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/nbio/st"
)

func TestMatchURL(t *testing.T) {
	defer gock.Disable()

	gock.New("http://(.*).com").
		Reply(200).
		BodyString("foo foo")

	res, err := http.Get("http://foo.com")
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 200)
	body, _ := io.ReadAll(res.Body)
	st.Expect(t, string(body), "foo foo")
}
