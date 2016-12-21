package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func TestCORSIsOriginAllowed(t *testing.T) {
	cors := &CORS{AllowOrigins: []string{"*"}}
	cors.generatePatterns()

	// When all the origins are accepted
	gobottest.Assert(t, cors.isOriginAllowed("http://localhost:8000"), true)
	gobottest.Assert(t, cors.isOriginAllowed("http://localhost:3001"), true)
	gobottest.Assert(t, cors.isOriginAllowed("http://server.com"), true)

	// When one origin is accepted
	cors = &CORS{AllowOrigins: []string{"http://localhost:8000"}}
	cors.generatePatterns()

	gobottest.Assert(t, cors.isOriginAllowed("http://localhost:8000"), true)
	gobottest.Assert(t, cors.isOriginAllowed("http://localhost:3001"), false)
	gobottest.Assert(t, cors.isOriginAllowed("http://server.com"), false)

	// When several origins are accepted
	cors = &CORS{AllowOrigins: []string{"http://localhost:*", "http://server.com"}}
	cors.generatePatterns()

	gobottest.Assert(t, cors.isOriginAllowed("http://localhost:8000"), true)
	gobottest.Assert(t, cors.isOriginAllowed("http://localhost:3001"), true)
	gobottest.Assert(t, cors.isOriginAllowed("http://server.com"), true)

	// When several origins are accepted within the same domain
	cors = &CORS{AllowOrigins: []string{"http://*.server.com"}}
	cors.generatePatterns()

	gobottest.Assert(t, cors.isOriginAllowed("http://localhost:8000"), false)
	gobottest.Assert(t, cors.isOriginAllowed("http://localhost:3001"), false)
	gobottest.Assert(t, cors.isOriginAllowed("http://foo.server.com"), true)
	gobottest.Assert(t, cors.isOriginAllowed("http://api.server.com"), true)
}

func TestCORSAllowedHeaders(t *testing.T) {
	cors := &CORS{AllowOrigins: []string{"*"}, AllowHeaders: []string{"Header1", "Header2"}}

	gobottest.Assert(t, cors.AllowedHeaders(), "Header1,Header2")
}

func TestCORSAllowedMethods(t *testing.T) {
	cors := &CORS{AllowOrigins: []string{"*"}, AllowMethods: []string{"GET", "POST"}}

	gobottest.Assert(t, cors.AllowedMethods(), "GET,POST")

	cors.AllowMethods = []string{"GET", "POST", "PUT"}

	gobottest.Assert(t, cors.AllowedMethods(), "GET,POST,PUT")
}

func TestCORS(t *testing.T) {
	api := initTestAPI()

	// Accepted origin
	allowedOrigin := []string{"http://server.com"}
	api.AddHandler(AllowRequestsFrom(allowedOrigin[0]))

	request, _ := http.NewRequest("GET", "/api/", nil)
	request.Header.Set("Origin", allowedOrigin[0])
	response := httptest.NewRecorder()
	api.ServeHTTP(response, request)
	gobottest.Assert(t, response.Header()["Access-Control-Allow-Origin"], allowedOrigin)

	// Not accepted Origin
	disallowedOrigin := []string{"http://disallowed.com"}
	request, _ = http.NewRequest("GET", "/api/", nil)
	request.Header.Set("Origin", disallowedOrigin[0])
	response = httptest.NewRecorder()
	api.ServeHTTP(response, request)
	gobottest.Refute(t, response.Header()["Access-Control-Allow-Origin"], disallowedOrigin)
	gobottest.Refute(t, response.Header()["Access-Control-Allow-Origin"], allowedOrigin)
}
