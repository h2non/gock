// This example shows how to enable the networking for a request to a local server
// and mock a second request to a remote server.
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"gopkg.in/h2non/gock.v1"
)

// Starts a local HTTP server in background
func startHTTPServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Local server received a GET request")

		res, err := http.Get("http://httpbin.org/nope")
		if err != nil {
			msg := fmt.Sprintf("Error from request to httpbin: %s", err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		body, _ := ioutil.ReadAll(res.Body)
		// MUST NOT get original body since the networking
		// wasn't enabled for this request
		fmt.Printf("Body From httpbin: %s\n", string(body))
		fmt.Printf("Status From httpbin: %s\n", res.Status)

		io.WriteString(w, "Local Response="+res.Header.Get("Server"))
	}))
}

func main() {
	defer gock.Disable()
	defer gock.DisableNetworking()

	srv := startHTTPServer()
	defer srv.Close()

	// Register our local server
	gock.New(srv.URL).
		EnableNetworking()

	gock.New("http://httpbin.org").
		Get("/nope").
		Reply(201).
		SetHeader("Server", "gock")

	res, err := http.Get(srv.URL)
	if err != nil {
		fmt.Printf("Error from request to localhost: %s", err)
		return
	}

	// The response status comes from the mock
	fmt.Printf("Status: %d\n", res.StatusCode)
	// The server header comes from mock as well
	fmt.Printf("Server header: %s\n", res.Header.Get("Server"))
	// MUST get original response since the networking was enabled for this request
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Printf("Body From Local Server: %s", string(body))
}
