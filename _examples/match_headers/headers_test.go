package test

import (
	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestMatchHeaders(t *testing.T) {
	defer gock.Disable()

	gock.New("http://foo.com").
		MatchHeader("Authorization", "^foo bar$").
		MatchHeader("API", "1.[0-9]+").
		HeaderPresent("Accept").
		Reply(200).
		BodyString("foo foo")

	req, err := http.NewRequest("GET", "http://foo.com", nil)
	req.Header.Set("Authorization", "foo bar")
	req.Header["API"] = []string{"1.0"}
	req.Header.Set("Accept", "text/plain")

	res, err := (&http.Client{}).Do(req)
	st.Expect(t, err, nil)
	st.Expect(t, res.StatusCode, 200)
	body, _ := ioutil.ReadAll(res.Body)
	st.Expect(t, string(body), "foo foo")
}
