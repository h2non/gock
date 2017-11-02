package test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/nbio/st"
)

func TestReplyError(t *testing.T) {
	defer gock.Off()

	gock.New("http://foo.com").
		Get("/bar").
		ReplyError(errors.New("Error dude!"))

	_, err := http.Get("http://foo.com/bar")
	st.Expect(t, err.Error(), "Get http://foo.com/bar: Error dude!")

	// Verify that we don't have pending mocks
	st.Expect(t, gock.IsDone(), true)
}
