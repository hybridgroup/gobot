package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hybridgroup/gobot"
)

func TestCORSIsOriginAllowed(t *testing.T) {
	cors := &CORS{AllowOrigins: []string{"*"}}
	cors.generatePatterns()

	// When all the origins are accepted
	gobot.Assert(t, cors.isOriginAllowed("http://localhost:8000"), true)
	gobot.Assert(t, cors.isOriginAllowed("http://localhost:3001"), true)
	gobot.Assert(t, cors.isOriginAllowed("http://server.com"), true)

	// When one origin is accepted
	cors = &CORS{AllowOrigins: []string{"http://localhost:8000"}}
	cors.generatePatterns()

	gobot.Assert(t, cors.isOriginAllowed("http://localhost:8000"), true)
	gobot.Assert(t, cors.isOriginAllowed("http://localhost:3001"), false)
	gobot.Assert(t, cors.isOriginAllowed("http://server.com"), false)

	// When several origins are accepted
	cors = &CORS{AllowOrigins: []string{"http://localhost:*", "http://server.com"}}
	cors.generatePatterns()

	gobot.Assert(t, cors.isOriginAllowed("http://localhost:8000"), true)
	gobot.Assert(t, cors.isOriginAllowed("http://localhost:3001"), true)
	gobot.Assert(t, cors.isOriginAllowed("http://server.com"), true)

	// When several origins are accepted within the same domain
	cors = &CORS{AllowOrigins: []string{"http://*.server.com"}}
	cors.generatePatterns()

	gobot.Assert(t, cors.isOriginAllowed("http://localhost:8000"), false)
	gobot.Assert(t, cors.isOriginAllowed("http://localhost:3001"), false)
	gobot.Assert(t, cors.isOriginAllowed("http://foo.server.com"), true)
	gobot.Assert(t, cors.isOriginAllowed("http://api.server.com"), true)
}

func TestCORSAllowedHeaders(t *testing.T) {
	cors := &CORS{AllowOrigins: []string{"*"}, AllowHeaders: []string{"Header1", "Header2"}}

	gobot.Assert(t, cors.AllowedHeaders(), "Header1,Header2")
}

func TestCORSAllowedMethods(t *testing.T) {
	cors := &CORS{AllowOrigins: []string{"*"}, AllowMethods: []string{"GET", "POST"}}

	gobot.Assert(t, cors.AllowedMethods(), "GET,POST")

	cors.AllowMethods = []string{"GET", "POST", "PUT"}

	gobot.Assert(t, cors.AllowedMethods(), "GET,POST,PUT")
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
	gobot.Assert(t, response.Header()["Access-Control-Allow-Origin"], allowedOrigin)

	// Not accepted Origin
	disallowedOrigin := []string{"http://disallowed.com"}
	request, _ = http.NewRequest("GET", "/api/", nil)
	request.Header.Set("Origin", disallowedOrigin[0])
	response = httptest.NewRecorder()
	api.ServeHTTP(response, request)
	gobot.Refute(t, response.Header()["Access-Control-Allow-Origin"], disallowedOrigin)
	gobot.Refute(t, response.Header()["Access-Control-Allow-Origin"], allowedOrigin)
}
