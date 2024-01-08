package gock

import (
	"net/http"

	"github.com/h2non/gock/threadsafe"
)

// Responder builds a mock http.Response based on the given Response mock.
func Responder(req *http.Request, mock *Response, res *http.Response) (*http.Response, error) {
	return threadsafe.Responder(req, mock, res)
}
