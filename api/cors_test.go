//nolint:usestdlibvars,noctx // ok here
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORSIsOriginAllowed(t *testing.T) {
	cors := &CORS{AllowOrigins: []string{"*"}}
	cors.generatePatterns()

	// When all the origins are accepted
	assert.True(t, cors.isOriginAllowed("http://localhost:8000"))
	assert.True(t, cors.isOriginAllowed("http://localhost:3001"))
	assert.True(t, cors.isOriginAllowed("http://server.com"))

	// When one origin is accepted
	cors = &CORS{AllowOrigins: []string{"http://localhost:8000"}}
	cors.generatePatterns()

	assert.True(t, cors.isOriginAllowed("http://localhost:8000"))
	assert.False(t, cors.isOriginAllowed("http://localhost:3001"))
	assert.False(t, cors.isOriginAllowed("http://server.com"))

	// When several origins are accepted
	cors = &CORS{AllowOrigins: []string{"http://localhost:*", "http://server.com"}}
	cors.generatePatterns()

	assert.True(t, cors.isOriginAllowed("http://localhost:8000"))
	assert.True(t, cors.isOriginAllowed("http://localhost:3001"))
	assert.True(t, cors.isOriginAllowed("http://server.com"))

	// When several origins are accepted within the same domain
	cors = &CORS{AllowOrigins: []string{"http://*.server.com"}}
	cors.generatePatterns()

	assert.False(t, cors.isOriginAllowed("http://localhost:8000"))
	assert.False(t, cors.isOriginAllowed("http://localhost:3001"))
	assert.True(t, cors.isOriginAllowed("http://foo.server.com"))
	assert.True(t, cors.isOriginAllowed("http://api.server.com"))
}

func TestCORSAllowedHeaders(t *testing.T) {
	cors := &CORS{AllowOrigins: []string{"*"}, AllowHeaders: []string{"Header1", "Header2"}}

	assert.Equal(t, "Header1,Header2", cors.AllowedHeaders())
}

func TestCORSAllowedMethods(t *testing.T) {
	cors := &CORS{AllowOrigins: []string{"*"}, AllowMethods: []string{"GET", "POST"}}

	assert.Equal(t, "GET,POST", cors.AllowedMethods())

	cors.AllowMethods = []string{"GET", "POST", "PUT"}

	assert.Equal(t, "GET,POST,PUT", cors.AllowedMethods())
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
	assert.Equal(t, allowedOrigin, response.Header()["Access-Control-Allow-Origin"])

	// Not accepted Origin
	disallowedOrigin := []string{"http://disallowed.com"}
	request, _ = http.NewRequest("GET", "/api/", nil)
	request.Header.Set("Origin", disallowedOrigin[0])
	response = httptest.NewRecorder()
	api.ServeHTTP(response, request)
	assert.NotEqual(t, disallowedOrigin, response.Header()["Access-Control-Allow-Origin"])
	assert.NotEqual(t, allowedOrigin, response.Header()["Access-Control-Allow-Origin"])
}
