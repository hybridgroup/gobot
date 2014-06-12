package api

import (
	"bytes"
	"encoding/json"
	"github.com/hybridgroup/gobot"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type null struct{}

func (null) Write(p []byte) (int, error) {
	return len(p), nil
}

var m *gobot.Gobot
var a *api

func init() {
	log.SetOutput(new(null))
	m = gobot.NewGobot()
	a = NewAPI(m)
	a.start = func(m *api) {}
	a.Start()

	m.Robots = []*gobot.Robot{
		gobot.NewTestRobot("Robot 1"),
		gobot.NewTestRobot("Robot 2"),
		gobot.NewTestRobot("Robot 3"),
	}
}

func TestRobots(t *testing.T) {
	request, _ := http.NewRequest("GET", "/robots", nil)
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i []map[string]interface{}
	json.Unmarshal(body, &i)
	gobot.Expect(t, len(i), 3)
}

func TestRobot(t *testing.T) {
	request, _ := http.NewRequest("GET", "/robots/Robot%201", nil)
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i map[string]interface{}
	json.Unmarshal(body, &i)
	gobot.Expect(t, i["name"].(string), "Robot 1")
}

func TestRobotDevices(t *testing.T) {
	request, _ := http.NewRequest("GET", "/robots/Robot%201/devices", nil)
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i []map[string]interface{}
	json.Unmarshal(body, &i)
	gobot.Expect(t, len(i), 3)
}

func TestRobotCommands(t *testing.T) {
	request, _ := http.NewRequest("GET", "/robots/Robot%201/commands", nil)
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i []string
	json.Unmarshal(body, &i)
	gobot.Expect(t, i, []string{"robotTestFunction"})
}
func TestExecuteRobotCommand(t *testing.T) {
	request, _ := http.NewRequest("GET", "/robots/Robot%201/commands/robotTestFunction", bytes.NewBufferString(`{"message":"Beep Boop"}`))
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i interface{}
	json.Unmarshal(body, &i)
	gobot.Expect(t, i, "hey Robot 1, Beep Boop")
}

func TestUnknownRobotCommand(t *testing.T) {
	request, _ := http.NewRequest("GET", "/robots/Robot%201/commands/robotTestFuntion1", bytes.NewBufferString(`{"message":"Beep Boop"}`))
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i interface{}
	json.Unmarshal(body, &i)
	gobot.Expect(t, i, "Unknown Command")
}

func TestRobotDevice(t *testing.T) {
	request, _ := http.NewRequest("GET", "/robots/Robot%201/devices/Device%201", nil)
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i map[string]interface{}
	json.Unmarshal(body, &i)
	gobot.Expect(t, i["name"].(string), "Device 1")
}

func TestRobotDeviceCommands(t *testing.T) {
	request, _ := http.NewRequest("GET", "/robots/Robot%201/devices/Device%201/commands", nil)
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i []string
	json.Unmarshal(body, &i)
	gobot.Expect(t, i, []string{"TestDriverCommand", "DriverCommand"})
}

func TestExecuteRobotDeviceCommand(t *testing.T) {
	request, _ := http.NewRequest("GET", "/robots/Robot%201/devices/Device%201/commands/TestDriverCommand", bytes.NewBufferString(`{"name":"human"}`))
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i interface{}
	json.Unmarshal(body, &i)
	gobot.Expect(t, i, "hello human")
}

func TestUnknownRobotDeviceCommand(t *testing.T) {
	request, _ := http.NewRequest("GET", "/robots/Robot%201/devices/Device%201/commands/DriverCommand1", bytes.NewBufferString(`{"name":"human"}`))
	request.Header.Add("Content-Type", "application/json")
	response := httptest.NewRecorder()
	a.server.ServeHTTP(response, request)

	body, _ := ioutil.ReadAll(response.Body)
	var i interface{}
	json.Unmarshal(body, &i)
	gobot.Expect(t, i, "Unknown Command")
}
