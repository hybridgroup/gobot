package particle

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/donovanhide/eventsource"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

// HELPERS

func createTestServer(handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

func getDummyResponseForPath(path string, dummy_response string, t *testing.T) *httptest.Server {
	dummy_data := []byte(dummy_response)

	return createTestServer(func(w http.ResponseWriter, r *http.Request) {
		actualPath := "/v1/devices" + path
		if r.URL.Path != actualPath {
			t.Errorf("Path doesn't match, expected %#v, got %#v", actualPath, r.URL.Path)
		}
		w.Write(dummy_data)
	})
}

func getDummyResponseForPathWithParams(path string, params []string, dummy_response string, t *testing.T) *httptest.Server {
	dummy_data := []byte(dummy_response)

	return createTestServer(func(w http.ResponseWriter, r *http.Request) {
		actualPath := "/v1/devices" + path
		if r.URL.Path != actualPath {
			t.Errorf("Path doesn't match, expected %#v, got %#v", actualPath, r.URL.Path)
		}

		r.ParseForm()

		for key, value := range params {
			if r.Form["params"][key] != value {
				t.Error("Expected param to be " + r.Form["params"][key] + " but was " + value)
			}
		}
		w.Write(dummy_data)
	})
}

func initTestAdaptor() *Adaptor {
	return NewAdaptor("myDevice", "token")
}

func initTestAdaptorWithServo() *Adaptor {
	a := NewAdaptor("myDevice", "token")
	a.servoPins["1"] = true
	return a
}

// TESTS

func TestAdaptor(t *testing.T) {
	var a interface{} = initTestAdaptor()
	_, ok := a.(gobot.Adaptor)
	if !ok {
		t.Errorf("Adaptor{} should be a gobot.Adaptor")
	}
}

func TestNewAdaptor(t *testing.T) {
	// does it return a pointer to an instance of Adaptor?
	var a interface{} = initTestAdaptor()
	core, ok := a.(*Adaptor)
	if !ok {
		t.Errorf("NewAdaptor() should have returned a *Adaptor")
	}

	gobottest.Assert(t, core.APIServer, "https://api.particle.io")
	gobottest.Assert(t, strings.HasPrefix(core.Name(), "Particle"), true)

	core.SetName("sparkie")
	gobottest.Assert(t, core.Name(), "sparkie")
}

func TestAdaptorConnect(t *testing.T) {
	a := initTestAdaptor()
	gobottest.Assert(t, a.Connect(), nil)
}

func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()
	a.Connect()
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestAdaptorAnalogRead(t *testing.T) {
	// When no error
	response := `{"return_value": 5.2}`
	params := []string{"A1"}

	a := initTestAdaptor()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/analogread", params, response, t)

	a.setAPIServer(testServer.URL)
	defer testServer.Close()

	val, _ := a.AnalogRead("A1")
	gobottest.Assert(t, val, 5)
}

func TestAdaptorAnalogReadError(t *testing.T) {
	a := initTestAdaptor()
	// When error
	testServer := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	defer testServer.Close()
	a.setAPIServer(testServer.URL)

	val, _ := a.AnalogRead("A1")
	gobottest.Assert(t, val, 0)
}

func TestAdaptorPwmWrite(t *testing.T) {
	response := `{}`
	params := []string{"A1,1"}

	a := initTestAdaptor()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/analogwrite", params, response, t)
	defer testServer.Close()

	a.setAPIServer(testServer.URL)
	a.PwmWrite("A1", 1)
}

func TestAdaptorAnalogWrite(t *testing.T) {
	response := `{}`
	params := []string{"A1,1"}

	a := initTestAdaptor()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/analogwrite", params, response, t)
	defer testServer.Close()

	a.setAPIServer(testServer.URL)
	a.AnalogWrite("A1", 1)
}

func TestAdaptorDigitalWrite(t *testing.T) {
	// When HIGH
	response := `{}`
	params := []string{"D7,HIGH"}

	a := initTestAdaptor()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/digitalwrite", params, response, t)

	a.setAPIServer(testServer.URL)
	a.DigitalWrite("D7", 1)

	testServer.Close()
	// When LOW
	params = []string{"D7,LOW"}

	testServer = getDummyResponseForPathWithParams("/"+a.DeviceID+"/digitalwrite", params, response, t)
	defer testServer.Close()

	a.setAPIServer(testServer.URL)
	a.DigitalWrite("D7", 0)
}

func TestAdaptorServoOpen(t *testing.T) {
	response := `{}`
	params := []string{"1"}

	a := initTestAdaptor()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/servoOpen", params, response, t)
	defer testServer.Close()

	a.setAPIServer(testServer.URL)
	a.servoPinOpen("1")
}

func TestAdaptorServoWrite(t *testing.T) {
	response := `{}`
	params := []string{"1,128"}

	a := initTestAdaptorWithServo()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/servoSet", params, response, t)
	defer testServer.Close()

	a.setAPIServer(testServer.URL)
	a.ServoWrite("1", 128)
}

func TestAdaptorDigitalRead(t *testing.T) {
	// When HIGH
	response := `{"return_value": 1}`
	params := []string{"D7"}

	a := initTestAdaptor()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/digitalread", params, response, t)

	a.setAPIServer(testServer.URL)

	val, _ := a.DigitalRead("D7")
	gobottest.Assert(t, val, 1)
	testServer.Close()

	// When LOW
	response = `{"return_value": 0}`

	testServer = getDummyResponseForPathWithParams("/"+a.DeviceID+"/digitalread", params, response, t)

	a.setAPIServer(testServer.URL)
	defer testServer.Close()

	val, _ = a.DigitalRead("D7")
	gobottest.Assert(t, val, 0)
}

func TestAdaptorDigitalReadError(t *testing.T) {
	a := initTestAdaptor()
	// When error
	testServer := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	defer testServer.Close()

	a.setAPIServer(testServer.URL)

	val, _ := a.DigitalRead("D7")
	gobottest.Assert(t, val, -1)
}

func TestAdaptorFunction(t *testing.T) {
	response := `{"return_value": 1}`

	a := initTestAdaptor()
	testServer := getDummyResponseForPath("/"+a.DeviceID+"/hello", response, t)

	a.setAPIServer(testServer.URL)

	val, _ := a.Function("hello", "100,200")
	gobottest.Assert(t, val, 1)
	testServer.Close()

	// When not existent
	response = `{"ok": false, "error": "timeout"}`
	testServer = getDummyResponseForPath("/"+a.DeviceID+"/hello", response, t)

	a.setAPIServer(testServer.URL)

	_, err := a.Function("hello", "")
	gobottest.Assert(t, err.Error(), "timeout")

	testServer.Close()
}

func TestAdaptorVariable(t *testing.T) {
	// When String
	response := `{"result": "1"}`

	a := initTestAdaptor()
	testServer := getDummyResponseForPath("/"+a.DeviceID+"/variable_name", response, t)

	a.setAPIServer(testServer.URL)

	val, _ := a.Variable("variable_name")
	gobottest.Assert(t, val, "1")
	testServer.Close()

	// When float
	response = `{"result": 1.1}`
	testServer = getDummyResponseForPath("/"+a.DeviceID+"/variable_name", response, t)

	a.setAPIServer(testServer.URL)

	val, _ = a.Variable("variable_name")
	gobottest.Assert(t, val, "1.1")
	testServer.Close()

	// When int
	response = `{"result": 1}`
	testServer = getDummyResponseForPath("/"+a.DeviceID+"/variable_name", response, t)

	a.setAPIServer(testServer.URL)

	val, _ = a.Variable("variable_name")
	gobottest.Assert(t, val, "1")
	testServer.Close()

	// When bool
	response = `{"result": true}`
	testServer = getDummyResponseForPath("/"+a.DeviceID+"/variable_name", response, t)

	a.setAPIServer(testServer.URL)

	val, _ = a.Variable("variable_name")
	gobottest.Assert(t, val, "true")
	testServer.Close()

	// When not existent
	response = `{"ok": false, "error": "Variable not found"}`
	testServer = getDummyResponseForPath("/"+a.DeviceID+"/not_existent", response, t)

	a.setAPIServer(testServer.URL)

	_, err := a.Variable("not_existent")
	gobottest.Assert(t, err.Error(), "Variable not found")

	testServer.Close()
}

func TestAdaptorSetAPIServer(t *testing.T) {
	a := initTestAdaptor()
	apiServer := "new_api_server"
	gobottest.Refute(t, a.APIServer, apiServer)

	a.setAPIServer(apiServer)
	gobottest.Assert(t, a.APIServer, apiServer)
}

func TestAdaptorDeviceURL(t *testing.T) {
	// When APIServer is set
	a := initTestAdaptor()
	a.setAPIServer("http://server")
	a.DeviceID = "devID"
	gobottest.Assert(t, a.deviceURL(), "http://server/v1/devices/devID")

	// When APIServer is not set
	a = &Adaptor{name: "particleie", DeviceID: "myDevice", AccessToken: "token"}
	gobottest.Assert(t, a.deviceURL(), "https://api.particle.io/v1/devices/myDevice")
}

func TestAdaptorPinLevel(t *testing.T) {
	a := initTestAdaptor()

	gobottest.Assert(t, a.pinLevel(1), "HIGH")
	gobottest.Assert(t, a.pinLevel(0), "LOW")
	gobottest.Assert(t, a.pinLevel(5), "LOW")
}

func TestAdaptorPostToparticle(t *testing.T) {
	a := initTestAdaptor()

	// When error on request
	vals := url.Values{}
	vals.Add("error", "error")
	resp, err := a.request("POST", "http://invalid%20host.com", vals)
	if err == nil {
		t.Error("request() should return an error when request was unsuccessful but returned", resp)
	}

	// When error reading body
	// Pending

	// When response.Status is not 200
	testServer := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	defer testServer.Close()

	resp, err = a.request("POST", testServer.URL+"/existent", vals)
	if err == nil {
		t.Error("request() should return an error when status is not 200 but returned", resp)
	}
}

type testEventSource struct {
	event string
	data  string
}

func (testEventSource) Id() string      { return "" }
func (t testEventSource) Event() string { return t.event }
func (t testEventSource) Data() string  { return t.data }

func TestAdaptorEventStream(t *testing.T) {
	a := initTestAdaptor()
	var url string
	eventSource = func(u string) (chan eventsource.Event, chan error, error) {
		url = u
		return nil, nil, nil
	}
	a.EventStream("all", "ping")
	gobottest.Assert(t, url, "https://api.particle.io/v1/events/ping?access_token=token")

	a.EventStream("devices", "ping")
	gobottest.Assert(t, url, "https://api.particle.io/v1/devices/events/ping?access_token=token")

	a.EventStream("device", "ping")
	gobottest.Assert(t, url, "https://api.particle.io/v1/devices/myDevice/events/ping?access_token=token")

	_, err := a.EventStream("nothing", "ping")
	gobottest.Assert(t, err.Error(), "source param should be: all, devices or device")

	eventSource = func(u string) (chan eventsource.Event, chan error, error) {
		return nil, nil, errors.New("error connecting sse")
	}

	_, err = a.EventStream("devices", "")
	gobottest.Assert(t, err.Error(), "error connecting sse")

	eventChan := make(chan eventsource.Event, 0)
	errorChan := make(chan error, 0)

	eventSource = func(u string) (chan eventsource.Event, chan error, error) {
		return eventChan, errorChan, nil
	}

	_, err = a.EventStream("devices", "")
	gobottest.Assert(t, err, nil)
}
