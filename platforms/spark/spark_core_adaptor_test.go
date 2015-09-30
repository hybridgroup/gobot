package spark

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/donovanhide/eventsource"
	"github.com/hybridgroup/gobot"
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

func initTestSparkCoreAdaptor() *SparkCoreAdaptor {
	return NewSparkCoreAdaptor("bot", "myDevice", "token")
}

// TESTS

func TestSparkCoreAdaptor(t *testing.T) {
	var _ gobot.Adaptor = (*SparkCoreAdaptor)(nil)

	var a interface{} = initTestSparkCoreAdaptor()
	_, ok := a.(gobot.Adaptor)
	if !ok {
		t.Errorf("SparkCoreAdaptor{} should be a gobot.Adaptor")
	}
}

func TestNewSparkCoreAdaptor(t *testing.T) {
	// does it return a pointer to an instance of SparkCoreAdaptor?
	var a interface{} = initTestSparkCoreAdaptor()
	spark, ok := a.(*SparkCoreAdaptor)
	if !ok {
		t.Errorf("NewSparkCoreAdaptor() should have returned a *SparkCoreAdaptor")
	}

	gobot.Assert(t, spark.APIServer, "https://api.spark.io")
	gobot.Assert(t, spark.Name(), "bot")
}

func TestSparkCoreAdaptorConnect(t *testing.T) {
	a := initTestSparkCoreAdaptor()
	gobot.Assert(t, len(a.Connect()), 0)
}

func TestSparkCoreAdaptorFinalize(t *testing.T) {
	a := initTestSparkCoreAdaptor()

	a.Connect()

	gobot.Assert(t, len(a.Finalize()), 0)
}

func TestSparkCoreAdaptorAnalogRead(t *testing.T) {
	// When no error
	response := `{"return_value": 5.2}`
	params := []string{"A1"}

	a := initTestSparkCoreAdaptor()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/analogread", params, response, t)

	a.setAPIServer(testServer.URL)
	defer testServer.Close()

	val, _ := a.AnalogRead("A1")
	gobot.Assert(t, val, 5)
}

func TestSparkCoreAdaptorAnalogReadError(t *testing.T) {
	a := initTestSparkCoreAdaptor()
	// When error
	testServer := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	defer testServer.Close()
	a.setAPIServer(testServer.URL)

	val, _ := a.AnalogRead("A1")
	gobot.Assert(t, val, 0)
}

func TestSparkCoreAdaptorPwmWrite(t *testing.T) {
	response := `{}`
	params := []string{"A1,1"}

	a := initTestSparkCoreAdaptor()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/analogwrite", params, response, t)
	defer testServer.Close()

	a.setAPIServer(testServer.URL)
	a.PwmWrite("A1", 1)
}

func TestSparkCoreAdaptorAnalogWrite(t *testing.T) {
	response := `{}`
	params := []string{"A1,1"}

	a := initTestSparkCoreAdaptor()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/analogwrite", params, response, t)
	defer testServer.Close()

	a.setAPIServer(testServer.URL)
	a.AnalogWrite("A1", 1)
}

func TestSparkCoreAdaptorDigitalWrite(t *testing.T) {
	// When HIGH
	response := `{}`
	params := []string{"D7,HIGH"}

	a := initTestSparkCoreAdaptor()
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

func TestSparkCoreAdaptorDigitalRead(t *testing.T) {
	// When HIGH
	response := `{"return_value": 1}`
	params := []string{"D7"}

	a := initTestSparkCoreAdaptor()
	testServer := getDummyResponseForPathWithParams("/"+a.DeviceID+"/digitalread", params, response, t)

	a.setAPIServer(testServer.URL)

	val, _ := a.DigitalRead("D7")
	gobot.Assert(t, val, 1)
	testServer.Close()

	// When LOW
	response = `{"return_value": 0}`

	testServer = getDummyResponseForPathWithParams("/"+a.DeviceID+"/digitalread", params, response, t)

	a.setAPIServer(testServer.URL)
	defer testServer.Close()

	val, _ = a.DigitalRead("D7")
	gobot.Assert(t, val, 0)
}

func TestSparkCoreAdaptorDigitalReadError(t *testing.T) {
	a := initTestSparkCoreAdaptor()
	// When error
	testServer := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	defer testServer.Close()

	a.setAPIServer(testServer.URL)

	val, _ := a.DigitalRead("D7")
	gobot.Assert(t, val, -1)
}

func TestSparkCoreAdaptorFunction(t *testing.T) {
	response := `{"return_value": 1}`

	a := initTestSparkCoreAdaptor()
	testServer := getDummyResponseForPath("/"+a.DeviceID+"/hello", response, t)

	a.setAPIServer(testServer.URL)

	val, _ := a.Function("hello", "100,200")
	gobot.Assert(t, val, 1)
	testServer.Close()

	// When not existent
	response = `{"ok": false, "error": "timeout"}`
	testServer = getDummyResponseForPath("/"+a.DeviceID+"/hello", response, t)

	a.setAPIServer(testServer.URL)

	_, err := a.Function("hello", "")
	gobot.Assert(t, err.Error(), "timeout")

	testServer.Close()
}

func TestSparkCoreAdaptorVariable(t *testing.T) {
	// When String
	response := `{"result": "1"}`

	a := initTestSparkCoreAdaptor()
	testServer := getDummyResponseForPath("/"+a.DeviceID+"/variable_name", response, t)

	a.setAPIServer(testServer.URL)

	val, _ := a.Variable("variable_name")
	gobot.Assert(t, val, "1")
	testServer.Close()

	// When float
	response = `{"result": 1.1}`
	testServer = getDummyResponseForPath("/"+a.DeviceID+"/variable_name", response, t)

	a.setAPIServer(testServer.URL)

	val, _ = a.Variable("variable_name")
	gobot.Assert(t, val, "1.1")
	testServer.Close()

	// When int
	response = `{"result": 1}`
	testServer = getDummyResponseForPath("/"+a.DeviceID+"/variable_name", response, t)

	a.setAPIServer(testServer.URL)

	val, _ = a.Variable("variable_name")
	gobot.Assert(t, val, "1")
	testServer.Close()

	// When bool
	response = `{"result": true}`
	testServer = getDummyResponseForPath("/"+a.DeviceID+"/variable_name", response, t)

	a.setAPIServer(testServer.URL)

	val, _ = a.Variable("variable_name")
	gobot.Assert(t, val, "true")
	testServer.Close()

	// When not existent
	response = `{"ok": false, "error": "Variable not found"}`
	testServer = getDummyResponseForPath("/"+a.DeviceID+"/not_existent", response, t)

	a.setAPIServer(testServer.URL)

	_, err := a.Variable("not_existent")
	gobot.Assert(t, err.Error(), "Variable not found")

	testServer.Close()
}

func TestSparkCoreAdaptorSetAPIServer(t *testing.T) {
	a := initTestSparkCoreAdaptor()
	apiServer := "new_api_server"
	gobot.Refute(t, a.APIServer, apiServer)

	a.setAPIServer(apiServer)
	gobot.Assert(t, a.APIServer, apiServer)
}

func TestSparkCoreAdaptorDeviceURL(t *testing.T) {
	// When APIServer is set
	a := initTestSparkCoreAdaptor()
	a.setAPIServer("http://server")
	a.DeviceID = "devID"
	gobot.Assert(t, a.deviceURL(), "http://server/v1/devices/devID")

	//When APIServer is not set
	a = &SparkCoreAdaptor{name: "sparkie", DeviceID: "myDevice", AccessToken: "token"}

	gobot.Assert(t, a.deviceURL(), "https://api.spark.io/v1/devices/myDevice")
}

func TestSparkCoreAdaptorPinLevel(t *testing.T) {

	a := initTestSparkCoreAdaptor()

	gobot.Assert(t, a.pinLevel(1), "HIGH")
	gobot.Assert(t, a.pinLevel(0), "LOW")
	gobot.Assert(t, a.pinLevel(5), "LOW")
}

func TestSparkCoreAdaptorPostToSpark(t *testing.T) {

	a := initTestSparkCoreAdaptor()

	// When error on request
	vals := url.Values{}
	vals.Add("error", "error")
	resp, err := a.requestToSpark("POST", "http://invalid%20host.com", vals)
	if err == nil {
		t.Errorf("requestToSpark() should return an error when request was unsuccessful but returned", resp)
	}

	// When error reading body
	// Pending

	// When response.Status is not 200
	testServer := createTestServer(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	defer testServer.Close()

	resp, err = a.requestToSpark("POST", testServer.URL+"/existent", vals)
	if err == nil {
		t.Errorf("requestToSpark() should return an error when status is not 200 but returned", resp)
	}

}

type testEventSource struct {
	event string
	data  string
}

func (testEventSource) Id() string      { return "" }
func (t testEventSource) Event() string { return t.event }
func (t testEventSource) Data() string  { return t.data }

func TestSparkCoreAdaptorEventStream(t *testing.T) {
	a := initTestSparkCoreAdaptor()
	var url string
	eventSource = func(u string) (chan eventsource.Event, chan error, error) {
		url = u
		return nil, nil, nil
	}
	a.EventStream("all", "ping")
	gobot.Assert(t, url, "https://api.spark.io/v1/events/ping?access_token=token")
	a.EventStream("devices", "ping")
	gobot.Assert(t, url, "https://api.spark.io/v1/devices/events/ping?access_token=token")
	a.EventStream("device", "ping")
	gobot.Assert(t, url, "https://api.spark.io/v1/devices/myDevice/events/ping?access_token=token")
	_, err := a.EventStream("nothing", "ping")
	gobot.Assert(t, err.Error(), "source param should be: all, devices or device")

	eventSource = func(u string) (chan eventsource.Event, chan error, error) {
		return nil, nil, errors.New("error connecting sse")
	}

	_, err = a.EventStream("devices", "")
	gobot.Assert(t, err.Error(), "error connecting sse")

	eventChan := make(chan eventsource.Event, 0)
	errorChan := make(chan error, 0)

	eventSource = func(u string) (chan eventsource.Event, chan error, error) {
		return eventChan, errorChan, nil
	}

	sem := make(chan bool, 0)
	stream, err := a.EventStream("devices", "")
	gobot.Assert(t, err, nil)

	// stream message
	gobot.Once(stream, func(data interface{}) {
		e := data.(Event)
		gobot.Assert(t, e.Name, "event")
		gobot.Assert(t, e.Data, "sse event")
		gobot.Assert(t, e.Error, nil)
		sem <- true
	})

	eventChan <- testEventSource{event: "event", data: "sse event"}

	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		t.Errorf("Did not recieve stream")
	}

	// stream error
	gobot.Once(stream, func(data interface{}) {
		e := data.(Event)
		gobot.Assert(t, e.Name, "")
		gobot.Assert(t, e.Data, "")
		gobot.Assert(t, e.Error.Error(), "stream error")
		sem <- true
	})

	errorChan <- errors.New("stream error")

	select {
	case <-sem:
	case <-time.After(1 * time.Second):
		t.Errorf("Did not recieve stream error")
	}

}
