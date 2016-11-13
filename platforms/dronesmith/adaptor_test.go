package dronesmith

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
)

var _ gobot.Adaptor = (*Adaptor)(nil)

// HELPERS

func createTestServer(handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

func getDummyResponseForPath(path string, dummy_response string, t *testing.T) *httptest.Server {
	dummy_data := []byte(dummy_response)

	return createTestServer(func(w http.ResponseWriter, r *http.Request) {
		actualPath := "/api/drone" + path
		if r.URL.Path != actualPath {
			t.Errorf("Path doesn't match, expected %#v, got %#v", actualPath, r.URL.Path)
		}
		w.Write(dummy_data)
	})
}

func initTestAdaptor() *Adaptor {
	return NewAdaptor("droneid", "email", "key")
}

// TESTS

func TestAdaptor(t *testing.T) {
	var _ gobot.Adaptor = (*Adaptor)(nil)

	var a interface{} = initTestAdaptor()
	_, ok := a.(gobot.Adaptor)
	if !ok {
		t.Errorf("Adaptor{} should be a gobot.Adaptor")
	}
}

func TestNewAdaptor(t *testing.T) {
	// does it return a pointer to an instance of Adaptor?
	var a interface{} = initTestAdaptor()
	drone, ok := a.(*Adaptor)
	if !ok {
		t.Errorf("NewAdaptor() should have returned a *Adaptor")
	}

	gobottest.Assert(t, drone.APIServer, "http://api.dronesmith.io")
	gobottest.Assert(t, drone.Name(), "Dronesmith")
}

// func TestAdaptorConnect(t *testing.T) {
// 	a := initTestAdaptor()
// 	gobottest.Assert(t, a.Connect(), nil)
// }
//
// func TestAdaptorFinalize(t *testing.T) {
// 	a := initTestAdaptor()
// 	a.Connect()
// 	gobottest.Assert(t, a.Finalize(), nil)
// }
