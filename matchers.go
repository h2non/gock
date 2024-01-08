package gock

import (
	"net/http"

	"github.com/h2non/gock/threadsafe"
)

// EOL represents the end of line character.
const EOL = threadsafe.EOL

// BodyTypes stores the supported MIME body types for matching.
// Currently only text-based types.
func BodyTypes() []string {
	return g.BodyTypes
}

func SetBodyTypes(types []string) {
	g.BodyTypes = types
}

// BodyTypeAliases stores a generic MIME type by alias.
func BodyTypeAliases() map[string]string {
	return g.BodyTypeAliases
}

func SetBodyTypeAliases(aliases map[string]string) {
	g.BodyTypeAliases = aliases
}

// CompressionSchemes stores the supported Content-Encoding types for decompression.
func CompressionSchemes() []string {
	return g.CompressionSchemes
}

func SetCompressionSchemes(schemes []string) {
	g.CompressionSchemes = schemes
}

// MatchMethod matches the HTTP method of the given request.
func MatchMethod(req *http.Request, ereq *Request) (bool, error) {
	return g.MatchMethod(req, ereq)
}

// MatchScheme matches the request URL protocol scheme.
func MatchScheme(req *http.Request, ereq *Request) (bool, error) {
	return g.MatchScheme(req, ereq)
}

// MatchHost matches the HTTP host header field of the given request.
func MatchHost(req *http.Request, ereq *Request) (bool, error) {
	return g.MatchHost(req, ereq)
}

// MatchPath matches the HTTP URL path of the given request.
func MatchPath(req *http.Request, ereq *Request) (bool, error) {
	return g.MatchPath(req, ereq)
}

// MatchHeaders matches the headers fields of the given request.
func MatchHeaders(req *http.Request, ereq *Request) (bool, error) {
	return g.MatchHeaders(req, ereq)
}

// MatchQueryParams matches the URL query params fields of the given request.
func MatchQueryParams(req *http.Request, ereq *Request) (bool, error) {
	return g.MatchQueryParams(req, ereq)
}

// MatchPathParams matches the URL path parameters of the given request.
func MatchPathParams(req *http.Request, ereq *Request) (bool, error) {
	return g.MatchPathParams(req, ereq)
}

// MatchBody tries to match the request body.
// TODO: not too smart now, needs several improvements.
func MatchBody(req *http.Request, ereq *Request) (bool, error) {
	return g.MatchBody(req, ereq)
}
