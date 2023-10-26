package particle

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/donovanhide/eventsource"
	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

// HELPERS

func createTestServer(handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

func getDummyResponseForPath(path string, dummyResponse string, t *testing.T) *httptest.Server {
	dummyData := []byte(dummyResponse)

	return createTestServer(func(w http.ResponseWriter, r *http.Request) {
		actualPath := "/v1/devices" + path
		if r.URL.Path != actualPath {
			t.Errorf("Path doesn't match, expected %#v, got %#v", actualPath, r.URL.Path)
		}
		_, _ = w.Write(dummyData)
	})
}

func getDummyResponseForPathWithParams(path string, params []string, dummyResponse string, t *testing.T) *httptest.Server {
	dummyData := []byte(dummyResponse)

	return createTestServer(func(w http.ResponseWriter, r *http.Request) {
		actualPath := "/v1/devices" + path
		if r.URL.Path != actualPath {
			t.Errorf("Path doesn't match, expected %#v, got %#v", actualPath, r.URL.Path)
		}

		_ = r.ParseForm()

		for key, value := range params {
			if r.Form["params"][key] != value {
				t.Error("Expected param to be " + r.Form["params"][key] + " but was " + value)
			}
		}
		_, _ = w.Write(dummyData)
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

	assert.Equal(t, "https://api.particle.io", core.APIServer)
	assert.True(t, strings.HasPrefix(core.Name(), "Particle"))

	core.SetName("sparkie")
	assert.Equal(t, "sparkie", core.Name())
}

func TestAdaptorConnect(t *testing.T) {
	a := initTestAdaptor()
	assert.NoError(t, a.Connect())
}

func TestAdaptorFinalize(t *testing.T) {
	a := initTestAdaptor()
	_ = a.Connect()
	assert.NoError(t, a.Finalize())
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
	assert.Equal(t, 5, val)
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
	assert.Equal(t, 0, val)
}

func TestAdaptorPwmWrite(t *testing.T) {
	response := `{}`
	params := []string{"A1,1"}

	a := initTestAdaptor()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/analogwrite", params, response, t)
	defer testServer.Close()

	a.setAPIServer(testServer.URL)
	_ = a.PwmWrite("A1", 1)
}

func TestAdaptorAnalogWrite(t *testing.T) {
	response := `{}`
	params := []string{"A1,1"}

	a := initTestAdaptor()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/analogwrite", params, response, t)
	defer testServer.Close()

	a.setAPIServer(testServer.URL)
	_ = a.AnalogWrite("A1", 1)
}

func TestAdaptorDigitalWrite(t *testing.T) {
	// When HIGH
	response := `{}`
	params := []string{"D7,HIGH"}

	a := initTestAdaptor()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/digitalwrite", params, response, t)

	a.setAPIServer(testServer.URL)
	_ = a.DigitalWrite("D7", 1)

	testServer.Close()
	// When LOW
	params = []string{"D7,LOW"}

	testServer = getDummyResponseForPathWithParams("/"+a.DeviceID+"/digitalwrite", params, response, t)
	defer testServer.Close()

	a.setAPIServer(testServer.URL)
	_ = a.DigitalWrite("D7", 0)
}

func TestAdaptorServoOpen(t *testing.T) {
	response := `{}`
	params := []string{"1"}

	a := initTestAdaptor()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/servoOpen", params, response, t)
	defer testServer.Close()

	a.setAPIServer(testServer.URL)
	_ = a.servoPinOpen("1")
}

func TestAdaptorServoWrite(t *testing.T) {
	response := `{}`
	params := []string{"1,128"}

	a := initTestAdaptorWithServo()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/servoSet", params, response, t)
	defer testServer.Close()

	a.setAPIServer(testServer.URL)
	_ = a.ServoWrite("1", 128)
}

func TestAdaptorDigitalRead(t *testing.T) {
	// When HIGH
	response := `{"return_value": 1}`
	params := []string{"D7"}

	a := initTestAdaptor()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/digitalread", params, response, t)

	a.setAPIServer(testServer.URL)

	val, _ := a.DigitalRead("D7")
	assert.Equal(t, 1, val)
	testServer.Close()

	// When LOW
	response = `{"return_value": 0}`

	testServer = getDummyResponseForPathWithParams("/"+a.DeviceID+"/digitalread", params, response, t)

	a.setAPIServer(testServer.URL)
	defer testServer.Close()

	val, _ = a.DigitalRead("D7")
	assert.Equal(t, 0, val)
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
	assert.Equal(t, -1, val)
}

func TestAdaptorFunction(t *testing.T) {
	response := `{"return_value": 1}`

	a := initTestAdaptor()
	testServer := getDummyResponseForPath("/"+a.DeviceID+"/hello", response, t)

	a.setAPIServer(testServer.URL)

	val, _ := a.Function("hello", "100,200")
	assert.Equal(t, 1, val)
	testServer.Close()

	// When not existent
	response = `{"ok": false, "error": "timeout"}`
	testServer = getDummyResponseForPath("/"+a.DeviceID+"/hello", response, t)

	a.setAPIServer(testServer.URL)

	_, err := a.Function("hello", "")
	assert.ErrorContains(t, err, "timeout")

	testServer.Close()
}

func TestAdaptorVariable(t *testing.T) {
	// When String
	response := `{"result": "1"}`

	a := initTestAdaptor()
	testServer := getDummyResponseForPath("/"+a.DeviceID+"/variable_name", response, t)

	a.setAPIServer(testServer.URL)

	val, _ := a.Variable("variable_name")
	assert.Equal(t, "1", val)
	testServer.Close()

	// When float
	response = `{"result": 1.1}`
	testServer = getDummyResponseForPath("/"+a.DeviceID+"/variable_name", response, t)

	a.setAPIServer(testServer.URL)

	val, _ = a.Variable("variable_name")
	assert.Equal(t, "1.1", val)
	testServer.Close()

	// When int
	response = `{"result": 1}`
	testServer = getDummyResponseForPath("/"+a.DeviceID+"/variable_name", response, t)

	a.setAPIServer(testServer.URL)

	val, _ = a.Variable("variable_name")
	assert.Equal(t, "1", val)
	testServer.Close()

	// When bool
	response = `{"result": true}`
	testServer = getDummyResponseForPath("/"+a.DeviceID+"/variable_name", response, t)

	a.setAPIServer(testServer.URL)

	val, _ = a.Variable("variable_name")
	assert.Equal(t, "true", val)
	testServer.Close()

	// When not existent
	response = `{"ok": false, "error": "Variable not found"}`
	testServer = getDummyResponseForPath("/"+a.DeviceID+"/not_existent", response, t)

	a.setAPIServer(testServer.URL)

	_, err := a.Variable("not_existent")
	assert.ErrorContains(t, err, "Variable not found")

	testServer.Close()
}

func TestAdaptorSetAPIServer(t *testing.T) {
	a := initTestAdaptor()
	apiServer := "new_api_server"
	assert.NotEqual(t, apiServer, a.APIServer)

	a.setAPIServer(apiServer)
	assert.Equal(t, apiServer, a.APIServer)
}

func TestAdaptorDeviceURL(t *testing.T) {
	// When APIServer is set
	a := initTestAdaptor()
	a.setAPIServer("http://server")
	a.DeviceID = "devID"
	assert.Equal(t, "http://server/v1/devices/devID", a.deviceURL())

	// When APIServer is not set
	a = &Adaptor{name: "particleie", DeviceID: "myDevice", AccessToken: "token"}
	assert.Equal(t, "https://api.particle.io/v1/devices/myDevice", a.deviceURL())
}

func TestAdaptorPinLevel(t *testing.T) {
	a := initTestAdaptor()

	assert.Equal(t, "HIGH", a.pinLevel(1))
	assert.Equal(t, "LOW", a.pinLevel(0))
	assert.Equal(t, "LOW", a.pinLevel(5))
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

func TestAdaptorEventStream(t *testing.T) {
	a := initTestAdaptor()
	var url string
	eventSource = func(u string) (chan eventsource.Event, chan error, error) {
		url = u
		return nil, nil, nil
	}
	_, _ = a.EventStream("all", "ping")
	assert.Equal(t, "https://api.particle.io/v1/events/ping?access_token=token", url)

	_, _ = a.EventStream("devices", "ping")
	assert.Equal(t, "https://api.particle.io/v1/devices/events/ping?access_token=token", url)

	_, _ = a.EventStream("device", "ping")
	assert.Equal(t, "https://api.particle.io/v1/devices/myDevice/events/ping?access_token=token", url)

	_, err := a.EventStream("nothing", "ping")
	assert.ErrorContains(t, err, "source param should be: all, devices or device")

	eventSource = func(u string) (chan eventsource.Event, chan error, error) {
		return nil, nil, errors.New("error connecting sse")
	}

	_, err = a.EventStream("devices", "")
	assert.ErrorContains(t, err, "error connecting sse")

	eventChan := make(chan eventsource.Event)
	errorChan := make(chan error)

	eventSource = func(u string) (chan eventsource.Event, chan error, error) {
		return eventChan, errorChan, nil
	}

	_, err = a.EventStream("devices", "")
	assert.NoError(t, err)
}
