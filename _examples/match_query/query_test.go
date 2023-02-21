package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/nbio/st"
)

func TestMatchQueryParams(t *testing.T) {
	defer gock.Disable()

	gock.New("http://foo.com").
		MatchParam("foo", "^bar$").
		MatchParam("bar", "baz").
		ParamPresent("baz").
		Reply(200).
		BodyString("foo foo")

	req, err := http.NewRequest("GET", "http://foo.com?foo=bar&bar=baz&baz=foo", nil)
	res, err := (&http.Client{}).Do(req)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 200)
	body, _ := io.ReadAll(res.Body)
	st.Expect(t, string(body), "foo foo")
}
