package gobot

import (
	"bytes"
	"encoding/json"
	"github.com/go-martini/martini"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
)

var _ = Describe("Master", func() {
	var (
		myMaster *Master
		a        *api
	)

	BeforeEach(func() {
		myMaster = GobotMaster()
		startApi = func(m *martini.ClassicMartini) {}
		a = Api(myMaster)
		myMaster.Robots = []*Robot{
			newTestRobot("Robot 1"),
			newTestRobot("Robot 2"),
			newTestRobot("Robot 3"),
		}
		trap = func(c chan os.Signal) {
			c <- os.Interrupt
		}
		myMaster.Start()
	})

	Context("when valid", func() {
		It("should return all robots", func() {
			request, _ := http.NewRequest("GET", "/robots", nil)
			response := httptest.NewRecorder()
			a.robots(response, request)
			body, _ := ioutil.ReadAll(response.Body)
			var i []map[string]interface{}
			json.Unmarshal(body, &i)
			Expect(len(i)).To(Equal(3))
		})
		It("should return robot", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201", nil)
			response := httptest.NewRecorder()
			a.robot("Robot 1", response, request)
			body, _ := ioutil.ReadAll(response.Body)
			var i map[string]interface{}
			json.Unmarshal(body, &i)
			Expect(i["name"].(string)).To(Equal("Robot 1"))
		})
		It("should return all robot devices", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201/devices", nil)
			response := httptest.NewRecorder()
			a.robot_devices("Robot 1", response, request)
			body, _ := ioutil.ReadAll(response.Body)
			var i []map[string]interface{}
			json.Unmarshal(body, &i)
			Expect(len(i)).To(Equal(3))
		})
		It("should return robot commands", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201/commands", nil)
			response := httptest.NewRecorder()
			a.robot_commands("Robot 1", response, request)
			body, _ := ioutil.ReadAll(response.Body)
			var i []string
			json.Unmarshal(body, &i)
			Expect(i).To(Equal([]string{"robotTestFunction"}))
		})
		It("should execute robot command", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201/commands/robotTestFuntion", bytes.NewBufferString(`{"message":"Beep Boop"}`))
			request.Header.Add("Content-Type", "application/json")
			response := httptest.NewRecorder()
			a.executeRobotCommand("Robot 1", "robotTestFunction", response, request)
			body, _ := ioutil.ReadAll(response.Body)
			var i map[string]interface{}
			json.Unmarshal(body, &i)
			Expect(i["result"]).To(Equal("hey Robot 1, Beep Boop"))
		})
		It("should not execute unknown robot command", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201/commands/robotTestFuntion1", bytes.NewBufferString(`{"message":"Beep Boop"}`))
			request.Header.Add("Content-Type", "application/json")
			response := httptest.NewRecorder()
			a.executeRobotCommand("Robot 1", "robotTestFunction1", response, request)
			body, _ := ioutil.ReadAll(response.Body)
			var i map[string]interface{}
			json.Unmarshal(body, &i)
			Expect(i["result"]).To(Equal("Unknown Command"))
		})
		It("should return robot device", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201/devices/Device%201", nil)
			response := httptest.NewRecorder()
			a.robot_device("Robot 1", "Device 1", response, request)
			body, _ := ioutil.ReadAll(response.Body)
			var i map[string]interface{}
			json.Unmarshal(body, &i)
			Expect(i["name"].(string)).To(Equal("Device 1"))
		})
		It("should return device commands", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201/devices/Device%201/commands", nil)
			response := httptest.NewRecorder()
			a.robot_device_commands("Robot 1", "Device 1", response, request)
			body, _ := ioutil.ReadAll(response.Body)
			var i []string
			json.Unmarshal(body, &i)
			Expect(i).To(Equal([]string{"TestDriverCommand", "DriverCommand"}))
		})
		It("should execute device command", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201/devices/Device%201/commands/TestDriverCommand", bytes.NewBufferString(`{"name":"human"}`))
			request.Header.Add("Content-Type", "application/json")
			response := httptest.NewRecorder()
			a.executeCommand("Robot 1", "Device 1", "TestDriverCommand", response, request)
			body, _ := ioutil.ReadAll(response.Body)
			var i map[string]interface{}
			json.Unmarshal(body, &i)
			Expect(i["result"]).To(Equal("hello human"))
		})
		It("should not execute unknown device command", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201/devices/Device%201/commands/DriverCommand1", bytes.NewBufferString(`{"name":"human"}`))
			request.Header.Add("Content-Type", "application/json")
			response := httptest.NewRecorder()
			a.executeCommand("Robot 1", "Device 1", "DriverCommand4", response, request)
			body, _ := ioutil.ReadAll(response.Body)
			var i map[string]interface{}
			json.Unmarshal(body, &i)
			Expect(i["result"]).To(Equal("Unknown Command"))
		})
	})
})
