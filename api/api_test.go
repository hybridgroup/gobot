//nolint:forcetypeassert,usestdlibvars,bodyclose,noctx // ok here
package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2"
)

func initTestAPI() *API {
	log.SetOutput(NullReadWriteCloser{})
	g := gobot.NewManager()
	a := NewAPI(g)
	a.start = func(m *API) {}
	a.Start()
	a.Debug()

	g.AddRobot(newTestRobot("Robot1"))
	g.AddRobot(newTestRobot("Robot2"))
	g.AddRobot(newTestRobot("Robot3"))
	g.AddCommand("TestFunction", func(params map[string]interface{}) interface{} {
		message := params["message"].(string)
		return fmt.Sprintf("hey %v", message)
	})

	return a
}

func TestStartWithoutDefaults(t *testing.T) {
	log.SetOutput(NullReadWriteCloser{})
	g := gobot.NewManager()
	a := NewAPI(g)
	a.start = func(m *API) {}

	a.Get("/", func(res http.ResponseWriter, req *http.Request) {})
	a.StartWithoutDefaults()

	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
}

func TestRobeaux(t *testing.T) {
	a := initTestAPI()
	// html assets
	request, _ := http.NewRequest("GET", "/index.html", nil)
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
	// js assets
	request, _ = http.NewRequest("GET", "/js/script.js", nil)
	response = httptest.NewRecorder()
	a.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
	// css assets
	request, _ = http.NewRequest("GET", "/css/application.css", nil)
	response = httptest.NewRecorder()
	a.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
	// unknown asset
	request, _ = http.NewRequest("GET", "/js/fake/file.js", nil)
	response = httptest.NewRecorder()
	a.ServeHTTP(response, request)
	assert.Equal(t, 404, response.Code)
}

func TestIndex(t *testing.T) {
	a := initTestAPI()
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	a.ServeHTTP(response, request)

	assert.Equal(t, http.StatusMovedPermanently, response.Code)
	assert.Equal(t, "/index.html", response.Header()["Location"][0])
}

func TestMcp(t *testing.T) {
	a := initTestAPI()
	request, _ := http.NewRequest("GET", "/api/", nil)
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)

	var body map[string]interface{}
	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.NotNil(t, body["MCP"].(map[string]interface{})["robots"])
	assert.NotNil(t, body["MCP"].(map[string]interface{})["commands"])
}

func TestMcpCommands(t *testing.T) {
	a := initTestAPI()
	request, _ := http.NewRequest("GET", "/api/commands", nil)
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)

	var body map[string]interface{}
	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, []interface{}{"TestFunction"}, body["commands"])
}

func TestExecuteMcpCommand(t *testing.T) {
	var body interface{}
	a := initTestAPI()

	// known command
	request, _ := http.NewRequest("GET",
		"/api/commands/TestFunction",
		bytes.NewBufferString(`{"message":"Beep Boop"}`),
	)
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)

	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "hey Beep Boop", body.(map[string]interface{})["result"])

	// unknown command
	request, _ = http.NewRequest("GET",
		"/api/commands/TestFuntion1",
		bytes.NewBufferString(`{"message":"Beep Boop"}`),
	)
	request.Header.Add("Content-Type", "application/json")
	response = httptest.NewRecorder()
	a.ServeHTTP(response, request)

	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "Unknown Command", body.(map[string]interface{})["error"])
}

func TestRobots(t *testing.T) {
	a := initTestAPI()
	request, _ := http.NewRequest("GET", "/api/robots", nil)
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)

	var body map[string]interface{}
	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Len(t, body["robots"].([]interface{}), 3)
}

func TestRobot(t *testing.T) {
	a := initTestAPI()

	// known robot
	request, _ := http.NewRequest("GET", "/api/robots/Robot1", nil)
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)

	var body map[string]interface{}
	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "Robot1", body["robot"].(map[string]interface{})["name"].(string))

	// unknown robot
	request, _ = http.NewRequest("GET", "/api/robots/UnknownRobot1", nil)
	a.ServeHTTP(response, request)

	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "No Robot found with the name UnknownRobot1", body["error"])
}

func TestRobotDevices(t *testing.T) {
	a := initTestAPI()

	// known robot
	request, _ := http.NewRequest("GET", "/api/robots/Robot1/devices", nil)
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)

	var body map[string]interface{}
	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Len(t, body["devices"].([]interface{}), 3)

	// unknown robot
	request, _ = http.NewRequest("GET", "/api/robots/UnknownRobot1/devices", nil)
	a.ServeHTTP(response, request)

	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "No Robot found with the name UnknownRobot1", body["error"])
}

func TestRobotCommands(t *testing.T) {
	a := initTestAPI()

	// known robot
	request, _ := http.NewRequest("GET", "/api/robots/Robot1/commands", nil)
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)

	var body map[string]interface{}
	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, []interface{}{"robotTestFunction"}, body["commands"])

	// unknown robot
	request, _ = http.NewRequest("GET", "/api/robots/UnknownRobot1/commands", nil)
	a.ServeHTTP(response, request)

	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "No Robot found with the name UnknownRobot1", body["error"])
}

func TestExecuteRobotCommand(t *testing.T) {
	var body interface{}
	a := initTestAPI()
	// known command
	request, _ := http.NewRequest("GET",
		"/api/robots/Robot1/commands/robotTestFunction",
		bytes.NewBufferString(`{"message":"Beep Boop", "robot":"Robot1"}`),
	)
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)

	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "hey Robot1, Beep Boop", body.(map[string]interface{})["result"])

	// unknown command
	request, _ = http.NewRequest("GET",
		"/api/robots/Robot1/commands/robotTestFuntion1",
		bytes.NewBufferString(`{"message":"Beep Boop"}`),
	)
	request.Header.Add("Content-Type", "application/json")
	response = httptest.NewRecorder()
	a.ServeHTTP(response, request)

	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "Unknown Command", body.(map[string]interface{})["error"])

	// uknown robot
	request, _ = http.NewRequest("GET",
		"/api/robots/UnknownRobot1/commands/robotTestFuntion1",
		bytes.NewBufferString(`{"message":"Beep Boop"}`),
	)
	request.Header.Add("Content-Type", "application/json")
	a.ServeHTTP(response, request)

	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "No Robot found with the name UnknownRobot1", body.(map[string]interface{})["error"])
}

func TestRobotDevice(t *testing.T) {
	a := initTestAPI()

	// known device
	request, _ := http.NewRequest("GET",
		"/api/robots/Robot1/devices/Device1",
		nil,
	)
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)

	var body map[string]interface{}
	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "Device1", body["device"].(map[string]interface{})["name"].(string))

	// unknown device
	request, _ = http.NewRequest("GET",
		"/api/robots/Robot1/devices/UnknownDevice1", nil)
	a.ServeHTTP(response, request)

	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "No Device found with the name UnknownDevice1", body["error"])
}

func TestRobotDeviceCommands(t *testing.T) {
	a := initTestAPI()

	// known device
	request, _ := http.NewRequest("GET",
		"/api/robots/Robot1/devices/Device1/commands",
		nil,
	)
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)

	var body map[string]interface{}
	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Len(t, body["commands"].([]interface{}), 2)

	// unknown device
	request, _ = http.NewRequest("GET",
		"/api/robots/Robot1/devices/UnknownDevice1/commands",
		nil,
	)
	a.ServeHTTP(response, request)
	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "No Device found with the name UnknownDevice1", body["error"])
}

func TestExecuteRobotDeviceCommand(t *testing.T) {
	var body interface{}
	a := initTestAPI()

	// known command
	request, _ := http.NewRequest("GET",
		"/api/robots/Robot1/devices/Device1/commands/TestDriverCommand",
		bytes.NewBufferString(`{"name":"human"}`),
	)
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)

	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "hello human", body.(map[string]interface{})["result"].(string))

	// unknown command
	request, _ = http.NewRequest("GET",
		"/api/robots/Robot1/devices/Device1/commands/DriverCommand1",
		bytes.NewBufferString(`{"name":"human"}`),
	)
	request.Header.Add("Content-Type", "application/json")
	response = httptest.NewRecorder()
	a.ServeHTTP(response, request)

	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "Unknown Command", body.(map[string]interface{})["error"])

	// unknown device
	request, _ = http.NewRequest("GET",
		"/api/robots/Robot1/devices/UnknownDevice1/commands/DriverCommand1",
		bytes.NewBufferString(`{"name":"human"}`),
	)
	request.Header.Add("Content-Type", "application/json")
	a.ServeHTTP(response, request)

	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "No Device found with the name UnknownDevice1", body.(map[string]interface{})["error"])
}

func TestRobotConnections(t *testing.T) {
	a := initTestAPI()

	// known robot
	request, _ := http.NewRequest("GET", "/api/robots/Robot1/connections", nil)
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)

	var body map[string]interface{}
	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Len(t, body["connections"].([]interface{}), 3)

	// unknown robot
	request, _ = http.NewRequest("GET", "/api/robots/UnknownRobot1/connections", nil)
	a.ServeHTTP(response, request)

	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "No Robot found with the name UnknownRobot1", body["error"])
}

func TestRobotConnection(t *testing.T) {
	a := initTestAPI()

	// known connection
	request, _ := http.NewRequest("GET",
		"/api/robots/Robot1/connections/Connection1",
		nil,
	)
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)

	var body map[string]interface{}
	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "Connection1", body["connection"].(map[string]interface{})["name"].(string))

	// unknown connection
	request, _ = http.NewRequest("GET",
		"/api/robots/Robot1/connections/UnknownConnection1",
		nil,
	)
	a.ServeHTTP(response, request)
	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "No Connection found with the name UnknownConnection1", body["error"])
}

func TestRobotDeviceEvent(t *testing.T) {
	a := initTestAPI()
	server := httptest.NewServer(a)
	defer server.Close()

	eventsURL := "/api/robots/Robot1/devices/Device1/events/"

	// known event
	respc := make(chan *http.Response, 1)
	go func() {
		resp, _ := http.Get(server.URL + eventsURL + "TestEvent")
		respc <- resp
	}()

	event := a.manager.Robot("Robot1").
		Device("Device1").(gobot.Eventer).
		Event("TestEvent")

	go func() {
		time.Sleep(time.Millisecond * 10) // wait some time, so select below is ready
		a.manager.Robot("Robot1").
			Device("Device1").(gobot.Eventer).Publish(event, "event-data")
	}()

	select {
	case resp := <-respc:
		reader := bufio.NewReader(resp.Body)
		data, _ := reader.ReadString('\n')
		assert.Equal(t, "data: \"event-data\"\n", data)
	case <-time.After(50 * time.Millisecond):
		t.Error("Not receiving data")
	}

	server.CloseClientConnections()

	// unknown event
	response, _ := http.Get(server.URL + eventsURL + "UnknownEvent")

	var body map[string]interface{}
	_ = json.NewDecoder(response.Body).Decode(&body)
	assert.Equal(t, "No Event found with the name UnknownEvent", body["error"])
}

func TestAPIRouter(t *testing.T) {
	a := initTestAPI()

	a.Head("/test", func(res http.ResponseWriter, req *http.Request) {})
	request, _ := http.NewRequest("HEAD", "/test", nil)
	response := httptest.NewRecorder()
	a.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)

	a.Get("/test", func(res http.ResponseWriter, req *http.Request) {})
	request, _ = http.NewRequest("GET", "/test", nil)
	response = httptest.NewRecorder()
	a.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)

	a.Post("/test", func(res http.ResponseWriter, req *http.Request) {})
	request, _ = http.NewRequest("POST", "/test", nil)
	response = httptest.NewRecorder()
	a.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)

	a.Put("/test", func(res http.ResponseWriter, req *http.Request) {})
	request, _ = http.NewRequest("PUT", "/test", nil)
	response = httptest.NewRecorder()
	a.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)

	a.Delete("/test", func(res http.ResponseWriter, req *http.Request) {})
	request, _ = http.NewRequest("DELETE", "/test", nil)
	response = httptest.NewRecorder()
	a.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)

	a.Options("/test", func(res http.ResponseWriter, req *http.Request) {})
	request, _ = http.NewRequest("OPTIONS", "/test", nil)
	response = httptest.NewRecorder()
	a.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code)
}
