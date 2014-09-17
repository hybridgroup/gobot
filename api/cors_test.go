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
	allowedOrigins := []string{"http://server:port"}

	cors := NewCORS(allowedOrigins)

	gobot.Assert(t, cors.AllowOrigins, allowedOrigins)
}

func TestCORSIsOriginAllowed(t *testing.T) {
	cors := NewCORS([]string{"*"})

	// When all the origins are accepted
	gobot.Assert(t, cors.isOriginAllowed("http://localhost:8000"), true)
	gobot.Assert(t, cors.isOriginAllowed("http://localhost:3001"), true)
	gobot.Assert(t, cors.isOriginAllowed("http://server.com"), true)

	// When one origin is accepted
	cors.AllowOrigins = []string{"http://localhost:8000"}

	gobot.Assert(t, cors.isOriginAllowed("http://localhost:8000"), true)
	gobot.Assert(t, cors.isOriginAllowed("http://localhost:3001"), false)
	gobot.Assert(t, cors.isOriginAllowed("http://server.com"), false)

	// When several origins are accepted
	cors.AllowOrigins = []string{"http://localhost:8000", "http://server.com"}

	gobot.Assert(t, cors.isOriginAllowed("http://localhost:8000"), true)
	gobot.Assert(t, cors.isOriginAllowed("http://localhost:3001"), false)
	gobot.Assert(t, cors.isOriginAllowed("http://server.com"), true)
}
