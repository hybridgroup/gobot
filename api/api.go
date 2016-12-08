package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/bmizerany/pat"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api/robeaux"
)

// API represents an API server
type API struct {
	master   *gobot.Master
	router   *pat.PatternServeMux
	Host     string
	Port     string
	Cert     string
	Key      string
	handlers []func(http.ResponseWriter, *http.Request)
	start    func(*API)
}

// NewAPI returns a new api instance
func NewAPI(m *gobot.Master) *API {
	return &API{
		master: m,
		router: pat.New(),
		Port:   "3000",
		start: func(a *API) {
			log.Println("Initializing API on " + a.Host + ":" + a.Port + "...")
			http.Handle("/", a)

			go func() {
				if a.Cert != "" && a.Key != "" {
					http.ListenAndServeTLS(a.Host+":"+a.Port, a.Cert, a.Key, nil)
				} else {
					log.Println("WARNING: API using insecure connection. " +
						"We recommend using an SSL certificate with Gobot.")
					http.ListenAndServe(a.Host+":"+a.Port, nil)
				}
			}()
		},
	}
}

// ServeHTTP calls api handlers and then serves request using api router
func (a *API) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, handler := range a.handlers {
		rec := httptest.NewRecorder()
		handler(rec, req)
		for k, v := range rec.Header() {
			res.Header()[k] = v
		}
		if rec.Code == http.StatusUnauthorized {
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
			return
		}
	}
	a.router.ServeHTTP(res, req)
}

// Post wraps api router Post call
func (a *API) Post(path string, f func(http.ResponseWriter, *http.Request)) {
	a.router.Post(path, http.HandlerFunc(f))
}

// Put wraps api router Put call
func (a *API) Put(path string, f func(http.ResponseWriter, *http.Request)) {
	a.router.Put(path, http.HandlerFunc(f))
}

// Delete wraps api router Delete call
func (a *API) Delete(path string, f func(http.ResponseWriter, *http.Request)) {
	a.router.Del(path, http.HandlerFunc(f))
}

// Options wraps api router Options call
func (a *API) Options(path string, f func(http.ResponseWriter, *http.Request)) {
	a.router.Options(path, http.HandlerFunc(f))
}

// Get wraps api router Get call
func (a *API) Get(path string, f func(http.ResponseWriter, *http.Request)) {
	a.router.Get(path, http.HandlerFunc(f))
}

// Head wraps api router Head call
func (a *API) Head(path string, f func(http.ResponseWriter, *http.Request)) {
	a.router.Head(path, http.HandlerFunc(f))
}

// AddHandler appends handler to api handlers
func (a *API) AddHandler(f func(http.ResponseWriter, *http.Request)) {
	a.handlers = append(a.handlers, f)
}

// Start initializes the api by setting up c3pio routes and robeaux
func (a *API) Start() {
	mcpCommandRoute := "/api/commands/:command"
	robotDeviceCommandRoute := "/api/robots/:robot/devices/:device/commands/:command"
	robotCommandRoute := "/api/robots/:robot/commands/:command"

	a.Get("/api/commands", a.mcpCommands)
	a.Get(mcpCommandRoute, a.executeMcpCommand)
	a.Post(mcpCommandRoute, a.executeMcpCommand)
	a.Get("/api/robots", a.robots)
	a.Get("/api/robots/:robot", a.robot)
	a.Get("/api/robots/:robot/commands", a.robotCommands)
	a.Get(robotCommandRoute, a.executeRobotCommand)
	a.Post(robotCommandRoute, a.executeRobotCommand)
	a.Get("/api/robots/:robot/devices", a.robotDevices)
	a.Get("/api/robots/:robot/devices/:device", a.robotDevice)
	a.Get("/api/robots/:robot/devices/:device/events/:event", a.robotDeviceEvent)
	a.Get("/api/robots/:robot/devices/:device/commands", a.robotDeviceCommands)
	a.Get(robotDeviceCommandRoute, a.executeRobotDeviceCommand)
	a.Post(robotDeviceCommandRoute, a.executeRobotDeviceCommand)
	a.Get("/api/robots/:robot/connections", a.robotConnections)
	a.Get("/api/robots/:robot/connections/:connection", a.robotConnection)
	a.Get("/api/", a.mcp)

	a.Get("/", func(res http.ResponseWriter, req *http.Request) {
		http.Redirect(res, req, "/index.html", http.StatusMovedPermanently)
	})
	a.Get("/index.html", a.robeaux)
	a.Get("/images/:a", a.robeaux)
	a.Get("/js/:a", a.robeaux)
	a.Get("/js/:a/", a.robeaux)
	a.Get("/js/:a/:b", a.robeaux)
	a.Get("/css/:a", a.robeaux)
	a.Get("/css/:a/", a.robeaux)
	a.Get("/css/:a/:b", a.robeaux)
	a.Get("/partials/:a", a.robeaux)

	a.start(a)
}

// robeaux returns handler for robeaux routes.
// Writes asset in response and sets correct header
func (a *API) robeaux(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	buf, err := robeaux.Asset(path[1:])
	if err != nil {
		http.Error(res, err.Error(), http.StatusNotFound)
		return
	}
	t := strings.Split(path, ".")
	if t[len(t)-1] == "js" {
		res.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	} else if t[len(t)-1] == "css" {
		res.Header().Set("Content-Type", "text/css; charset=utf-8")
	} else if t[len(t)-1] == "html" {
		res.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	res.Write(buf)
}

// mcp returns MCP route handler.
// Writes JSON with gobot representation
func (a *API) mcp(res http.ResponseWriter, req *http.Request) {
	a.writeJSON(map[string]interface{}{"MCP": gobot.NewJSONMaster(a.master)}, res)
}

// mcpCommands returns commands route handler.
// Writes JSON with global commands representation
func (a *API) mcpCommands(res http.ResponseWriter, req *http.Request) {
	a.writeJSON(map[string]interface{}{"commands": gobot.NewJSONMaster(a.master).Commands}, res)
}

// robots returns route handler.
// Writes JSON with robots representation
func (a *API) robots(res http.ResponseWriter, req *http.Request) {
	jsonRobots := []*gobot.JSONRobot{}
	a.master.Robots().Each(func(r *gobot.Robot) {
		jsonRobots = append(jsonRobots, gobot.NewJSONRobot(r))
	})
	a.writeJSON(map[string]interface{}{"robots": jsonRobots}, res)
}

// robot returns route handler.
// Writes JSON with robot representation
func (a *API) robot(res http.ResponseWriter, req *http.Request) {
	if robot, err := a.jsonRobotFor(req.URL.Query().Get(":robot")); err != nil {
		a.writeJSON(map[string]interface{}{"error": err.Error()}, res)
	} else {
		a.writeJSON(map[string]interface{}{"robot": robot}, res)
	}
}

// robotCommands returns commands route handler
// Writes JSON with robot commands representation
func (a *API) robotCommands(res http.ResponseWriter, req *http.Request) {
	if robot, err := a.jsonRobotFor(req.URL.Query().Get(":robot")); err != nil {
		a.writeJSON(map[string]interface{}{"error": err.Error()}, res)
	} else {
		a.writeJSON(map[string]interface{}{"commands": robot.Commands}, res)
	}
}

// robotDevices returns devices route handler.
// Writes JSON with robot devices representation
func (a *API) robotDevices(res http.ResponseWriter, req *http.Request) {
	if robot := a.master.Robot(req.URL.Query().Get(":robot")); robot != nil {
		jsonDevices := []*gobot.JSONDevice{}
		robot.Devices().Each(func(d gobot.Device) {
			jsonDevices = append(jsonDevices, gobot.NewJSONDevice(d))
		})
		a.writeJSON(map[string]interface{}{"devices": jsonDevices}, res)
	} else {
		a.writeJSON(map[string]interface{}{"error": "No Robot found with the name " + req.URL.Query().Get(":robot")}, res)
	}
}

// robotDevice returns device route handler.
// Writes JSON with robot device representation
func (a *API) robotDevice(res http.ResponseWriter, req *http.Request) {
	if device, err := a.jsonDeviceFor(req.URL.Query().Get(":robot"), req.URL.Query().Get(":device")); err != nil {
		a.writeJSON(map[string]interface{}{"error": err.Error()}, res)
	} else {
		a.writeJSON(map[string]interface{}{"device": device}, res)
	}
}

func (a *API) robotDeviceEvent(res http.ResponseWriter, req *http.Request) {
	f, _ := res.(http.Flusher)
	c, _ := res.(http.CloseNotifier)

	dataChan := make(chan string)
	closer := c.CloseNotify()

	res.Header().Set("Content-Type", "text/event-stream")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")

	device := a.master.Robot(req.URL.Query().Get(":robot")).
		Device(req.URL.Query().Get(":device"))

	if event := a.master.Robot(req.URL.Query().Get(":robot")).
		Device(req.URL.Query().Get(":device")).(gobot.Eventer).
		Event(req.URL.Query().Get(":event")); len(event) > 0 {
		device.(gobot.Eventer).On(event, func(data interface{}) {
			d, _ := json.Marshal(data)
			dataChan <- string(d)
		})

		for {
			select {
			case data := <-dataChan:
				fmt.Fprintf(res, "data: %v\n\n", data)
				f.Flush()
			case <-closer:
				log.Println("Closing connection")
				return
			}
		}
	} else {
		a.writeJSON(map[string]interface{}{
			"error": "No Event found with the name " + req.URL.Query().Get(":event"),
		}, res)
	}
}

// robotDeviceCommands returns device commands route handler
// writes JSON with robot device commands representation
func (a *API) robotDeviceCommands(res http.ResponseWriter, req *http.Request) {
	if device, err := a.jsonDeviceFor(req.URL.Query().Get(":robot"), req.URL.Query().Get(":device")); err != nil {
		a.writeJSON(map[string]interface{}{"error": err.Error()}, res)
	} else {
		a.writeJSON(map[string]interface{}{"commands": device.Commands}, res)
	}
}

// robotConnections returns connections route handler
// writes JSON with robot connections representation
func (a *API) robotConnections(res http.ResponseWriter, req *http.Request) {
	jsonConnections := []*gobot.JSONConnection{}
	if robot := a.master.Robot(req.URL.Query().Get(":robot")); robot != nil {
		robot.Connections().Each(func(c gobot.Connection) {
			jsonConnections = append(jsonConnections, gobot.NewJSONConnection(c))
		})
		a.writeJSON(map[string]interface{}{"connections": jsonConnections}, res)
	} else {
		a.writeJSON(map[string]interface{}{"error": "No Robot found with the name " + req.URL.Query().Get(":robot")}, res)
	}

}

// robotConnection returns connection route handler
// writes JSON with robot connection representation
func (a *API) robotConnection(res http.ResponseWriter, req *http.Request) {
	if conn, err := a.jsonConnectionFor(req.URL.Query().Get(":robot"), req.URL.Query().Get(":connection")); err != nil {
		a.writeJSON(map[string]interface{}{"error": err.Error()}, res)
	} else {
		a.writeJSON(map[string]interface{}{"connection": conn}, res)
	}
}

// executeMcpCommand calls a global command associated to requested route
func (a *API) executeMcpCommand(res http.ResponseWriter, req *http.Request) {
	a.executeCommand(a.master.Command(req.URL.Query().Get(":command")),
		res,
		req,
	)
}

// executeRobotDeviceCommand calls a device command associated to requested route
func (a *API) executeRobotDeviceCommand(res http.ResponseWriter, req *http.Request) {
	if _, err := a.jsonDeviceFor(req.URL.Query().Get(":robot"),
		req.URL.Query().Get(":device")); err != nil {
		a.writeJSON(map[string]interface{}{"error": err.Error()}, res)
	} else {
		a.executeCommand(
			a.master.Robot(req.URL.Query().Get(":robot")).
				Device(req.URL.Query().Get(":device")).(gobot.Commander).
				Command(req.URL.Query().Get(":command")),
			res,
			req,
		)
	}
}

// executeRobotCommand calls a robot command associated to requested route
func (a *API) executeRobotCommand(res http.ResponseWriter, req *http.Request) {
	if _, err := a.jsonRobotFor(req.URL.Query().Get(":robot")); err != nil {
		a.writeJSON(map[string]interface{}{"error": err.Error()}, res)
	} else {
		a.executeCommand(
			a.master.Robot(req.URL.Query().Get(":robot")).
				Command(req.URL.Query().Get(":command")),
			res,
			req,
		)
	}
}

// executeCommand writes JSON response with `f` returned value.
func (a *API) executeCommand(f func(map[string]interface{}) interface{},
	res http.ResponseWriter,
	req *http.Request,
) {

	body := make(map[string]interface{})
	json.NewDecoder(req.Body).Decode(&body)

	if f != nil {
		a.writeJSON(map[string]interface{}{"result": f(body)}, res)
	} else {
		a.writeJSON(map[string]interface{}{"error": "Unknown Command"}, res)
	}
}

// writeJSON writes `j` as JSON in response
func (a *API) writeJSON(j interface{}, res http.ResponseWriter) {
	data, _ := json.Marshal(j)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

// Debug add handler to api that prints each request
func (a *API) Debug() {
	a.AddHandler(func(res http.ResponseWriter, req *http.Request) {
		log.Println(req)
	})
}

func (a *API) jsonRobotFor(name string) (jrobot *gobot.JSONRobot, err error) {
	if robot := a.master.Robot(name); robot != nil {
		jrobot = gobot.NewJSONRobot(robot)
	} else {
		err = errors.New("No Robot found with the name " + name)
	}
	return
}

func (a *API) jsonDeviceFor(robot string, name string) (jdevice *gobot.JSONDevice, err error) {
	if device := a.master.Robot(robot).Device(name); device != nil {
		jdevice = gobot.NewJSONDevice(device)
	} else {
		err = errors.New("No Device found with the name " + name)
	}
	return
}

func (a *API) jsonConnectionFor(robot string, name string) (jconnection *gobot.JSONConnection, err error) {
	if connection := a.master.Robot(robot).Connection(name); connection != nil {
		jconnection = gobot.NewJSONConnection(connection)
	} else {
		err = errors.New("No Connection found with the name " + name)
	}
	return
}
