package gock

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// Responder builds a mock http.Response based on the given Response mock.
func Responder(req *http.Request, mock *Response, res *http.Response) (*http.Response, error) {
	// If error present, reply it
	err := mock.Error
	if err != nil {
		return nil, err
	}

	if res == nil {
		res = createResponse(req)
	}

	// Apply response filter
	for _, filter := range mock.Filters {
		if !filter(res) {
			return res, nil
		}
	}

	// Define mock status code
	if mock.StatusCode != 0 {
		res.Status = strconv.Itoa(mock.StatusCode) + " " + http.StatusText(mock.StatusCode)
		res.StatusCode = mock.StatusCode
	}

	// Define headers by merging fields
	res.Header = mergeHeaders(res, mock)

	// Define mock body, if present
	if len(mock.BodyBuffer) > 0 {
		res.ContentLength = int64(len(mock.BodyBuffer))
		res.Body = createReadCloser(mock.BodyBuffer)
	}

	// Set raw mock body, if exist
	if mock.BodyGen != nil {
		res.ContentLength = -1
		res.Body = mock.BodyGen()
	}

	// Apply response mappers
	for _, mapper := range mock.Mappers {
		if tres := mapper(res); tres != nil {
			res = tres
		}
	}

	// Sleep to simulate delay, if necessary
	if mock.ResponseDelay > 0 {
		// allow escaping from sleep due to request context expiration or cancellation
		t := time.NewTimer(mock.ResponseDelay)
		select {
		case <-t.C:
		case <-req.Context().Done():
			// cleanly stop the timer
			if !t.Stop() {
				<-t.C
			}
		}
	}

	// check if the request context has ended. we could put this up in the delay code above, but putting it here
	// has the added benefit of working even when there is no delay (very small timeouts, already-done contexts, etc.)
	if err = req.Context().Err(); err != nil {
		// cleanly close the response and return the context error
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
		return nil, err
	}

	return res, err
}

// createResponse creates a new http.Response with default fields.
func createResponse(req *http.Request) *http.Response {
	return &http.Response{
		ProtoMajor: 1,
		ProtoMinor: 1,
		Proto:      "HTTP/1.1",
		Request:    req,
		Header:     make(http.Header),
		Body:       createReadCloser([]byte{}),
	}
}

// mergeHeaders copies the mock headers.
func mergeHeaders(res *http.Response, mres *Response) http.Header {
	for key, values := range mres.Header {
		for _, value := range values {
			res.Header.Add(key, value)
		}
	}
	return res.Header
}

// createReadCloser creates an io.ReadCloser from a byte slice that is suitable for use as an
// http response body.
func createReadCloser(body []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(body))
}
