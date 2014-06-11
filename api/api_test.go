package api

import (
	"bytes"
	"encoding/json"
	"github.com/hybridgroup/gobot"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("API", func() {
	var (
		m *gobot.Gobot
		a *api
	)

	BeforeEach(func() {
		m = gobot.NewGobot()
		a = NewAPI(m)
		a.start = func(m *api) {}
		a.Start()

		m.Robots = []*gobot.Robot{
			gobot.NewTestRobot("Robot 1"),
			gobot.NewTestRobot("Robot 2"),
			gobot.NewTestRobot("Robot 3"),
		}
	})

	Context("when valid", func() {
		It("should return all robots", func() {
			request, _ := http.NewRequest("GET", "/robots", nil)
			response := httptest.NewRecorder()
			a.server.ServeHTTP(response, request)

			body, _ := ioutil.ReadAll(response.Body)
			var i []map[string]interface{}
			json.Unmarshal(body, &i)
			Expect(len(i)).To(Equal(3))
		})
		It("should return robot", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201", nil)
			response := httptest.NewRecorder()
			a.server.ServeHTTP(response, request)

			body, _ := ioutil.ReadAll(response.Body)
			var i map[string]interface{}
			json.Unmarshal(body, &i)
			Expect(i["name"].(string)).To(Equal("Robot 1"))
		})
		It("should return all robot devices", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201/devices", nil)
			response := httptest.NewRecorder()
			a.server.ServeHTTP(response, request)

			body, _ := ioutil.ReadAll(response.Body)
			var i []map[string]interface{}
			json.Unmarshal(body, &i)
			Expect(len(i)).To(Equal(3))
		})
		PIt("should return robot commands", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201/commands", nil)
			response := httptest.NewRecorder()
			a.server.ServeHTTP(response, request)

			body, _ := ioutil.ReadAll(response.Body)
			var i []string
			json.Unmarshal(body, &i)
			Expect(i).To(Equal([]string{"robotTestFunction"}))
		})
		PIt("should execute robot command", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201/commands/robotTestFuntion", bytes.NewBufferString(`{"message":"Beep Boop"}`))
			request.Header.Add("Content-Type", "application/json")
			response := httptest.NewRecorder()
			a.server.ServeHTTP(response, request)

			body, _ := ioutil.ReadAll(response.Body)
			var i []interface{}
			json.Unmarshal(body, &i)
			Expect(i[0]).To(Equal("hey Robot 1, Beep Boop"))
		})
		It("should not execute unknown robot command", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201/commands/robotTestFuntion1", bytes.NewBufferString(`{"message":"Beep Boop"}`))
			request.Header.Add("Content-Type", "application/json")
			response := httptest.NewRecorder()
			a.server.ServeHTTP(response, request)

			body, _ := ioutil.ReadAll(response.Body)
			var i []interface{}
			json.Unmarshal(body, &i)
			Expect(i[0]).To(Equal("Unknown Command"))
		})
		It("should return robot device", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201/devices/Device%201", nil)
			response := httptest.NewRecorder()
			a.server.ServeHTTP(response, request)

			body, _ := ioutil.ReadAll(response.Body)
			var i map[string]interface{}
			json.Unmarshal(body, &i)
			Expect(i["name"].(string)).To(Equal("Device 1"))
		})
		It("should return device commands", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201/devices/Device%201/commands", nil)
			response := httptest.NewRecorder()
			a.server.ServeHTTP(response, request)

			body, _ := ioutil.ReadAll(response.Body)
			var i []string
			json.Unmarshal(body, &i)
			Expect(i).To(Equal([]string{"TestDriverCommand", "DriverCommand"}))
		})
		It("should execute device command", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201/devices/Device%201/commands/TestDriverCommand", bytes.NewBufferString(`{"name":"human"}`))
			request.Header.Add("Content-Type", "application/json")
			response := httptest.NewRecorder()
			a.server.ServeHTTP(response, request)

			body, _ := ioutil.ReadAll(response.Body)
			var i []interface{}
			json.Unmarshal(body, &i)
			Expect(i[0]).To(Equal("hello human"))
		})
		It("should not execute unknown device command", func() {
			request, _ := http.NewRequest("GET", "/robots/Robot%201/devices/Device%201/commands/DriverCommand1", bytes.NewBufferString(`{"name":"human"}`))
			request.Header.Add("Content-Type", "application/json")
			response := httptest.NewRecorder()
			a.server.ServeHTTP(response, request)

			body, _ := ioutil.ReadAll(response.Body)
			var i []interface{}
			json.Unmarshal(body, &i)
			Expect(i[0]).To(Equal("Unknown Command"))
		})
	})
})
