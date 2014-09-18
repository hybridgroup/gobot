package api

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

func TestNewCORS(t *testing.T) {
	var cors interface{} = NewCORS([]string{})

	// Does it return a pointer to an instance of CORS?
	_, ok := cors.(*CORS)
	if !ok {
		t.Errorf("NewCORS() should have returned a *CORS")
	}
}

func TestNewCorsSetsProperties(t *testing.T) {
	allowedOrigins := []string{"http://*server:*", "http://localhost:*"}
	allowedMethods := []string{"GET", "POST"}
	allowedHeaders := []string{"Origin", "Content-Type"}
	contentType := "application/json; charset=utf-8"
	allowOriginPatterns := []string{"^http://.*server:.*$", "^http://localhost:.*$"}

	cors := NewCORS(allowedOrigins)

	gobot.Assert(t, cors.AllowOrigins, allowedOrigins)
	gobot.Assert(t, cors.AllowMethods, allowedMethods)
	gobot.Assert(t, cors.AllowHeaders, allowedHeaders)
	gobot.Assert(t, cors.ContentType, contentType)
	gobot.Assert(t, cors.allowOriginPatterns, allowOriginPatterns)
}

func TestCORSIsOriginAllowed(t *testing.T) {
	cors := NewCORS([]string{"*"})

	// When all the origins are accepted
	gobot.Assert(t, cors.isOriginAllowed("http://localhost:8000"), true)
	gobot.Assert(t, cors.isOriginAllowed("http://localhost:3001"), true)
	gobot.Assert(t, cors.isOriginAllowed("http://server.com"), true)

	// When one origin is accepted
	cors = NewCORS([]string{"http://localhost:8000"})

	gobot.Assert(t, cors.isOriginAllowed("http://localhost:8000"), true)
	gobot.Assert(t, cors.isOriginAllowed("http://localhost:3001"), false)
	gobot.Assert(t, cors.isOriginAllowed("http://server.com"), false)

	// When several origins are accepted
	cors = NewCORS([]string{"http://localhost:*", "http://server.com"})

	gobot.Assert(t, cors.isOriginAllowed("http://localhost:8000"), true)
	gobot.Assert(t, cors.isOriginAllowed("http://localhost:3001"), true)
	gobot.Assert(t, cors.isOriginAllowed("http://server.com"), true)

	// When several origins are accepted within the same domain
	cors = NewCORS([]string{"http://*.server.com"})

	gobot.Assert(t, cors.isOriginAllowed("http://localhost:8000"), false)
	gobot.Assert(t, cors.isOriginAllowed("http://localhost:3001"), false)
	gobot.Assert(t, cors.isOriginAllowed("http://foo.server.com"), true)
	gobot.Assert(t, cors.isOriginAllowed("http://api.server.com"), true)
}

func TestCORSAllowedHeaders(t *testing.T) {
	cors := NewCORS([]string{"*"})

	cors.AllowHeaders = []string{"Header1", "Header2"}

	gobot.Assert(t, cors.AllowedHeaders(), "Header1,Header2")
}

func TestCORSAllowedMethods(t *testing.T) {
	cors := NewCORS([]string{"*"})

	gobot.Assert(t, cors.AllowedMethods(), "GET,POST")

	cors.AllowMethods = []string{"GET", "POST", "PUT"}

	gobot.Assert(t, cors.AllowedMethods(), "GET,POST,PUT")
}
