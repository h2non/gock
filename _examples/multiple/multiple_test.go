package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/nbio/st"
)

func TestMultipleMocks(t *testing.T) {
	defer gock.Disable()

	gock.New("http://server.com").
		Get("/foo").
		Reply(200).
		JSON(map[string]string{"value": "foo"})

	gock.New("http://server.com").
		Get("/bar").
		Reply(200).
		JSON(map[string]string{"value": "bar"})

	gock.New("http://server.com").
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

	for _, test := range tests {
		res, err := http.Get("http://server.com" + test.path)
		st.Expect(t, err, nil)
		st.Expect(t, res.StatusCode, 200)
		body, _ := io.ReadAll(res.Body)
		st.Expect(t, string(body)[:15], `{"value":"`+test.path[1:]+`"}`)
	}

	// Failed request after mocks expires
	_, err := http.Get("http://server.com/foo")
	st.Reject(t, err, nil)
}
