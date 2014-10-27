package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/bmizerany/pat"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api/robeaux"
)

type api struct {
	gobot    *gobot.Gobot
	router   *pat.PatternServeMux
	Host     string
	Port     string
	Cert     string
	Key      string
	handlers []func(http.ResponseWriter, *http.Request)
	start    func(*api)
}

// newAPI returns a gobot api instance
// and starts a http server using configuration options
func NewAPI(g *gobot.Gobot) *api {
	return &api{
		gobot:  g,
		router: pat.New(),
		Port:   "3000",
		start: func(a *api) {
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
func (a *api) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	for _, handler := range a.handlers {
		handler(res, req)
	}
	a.router.ServeHTTP(res, req)
}

// Post wraps api router Post call
func (a *api) Post(path string, f func(http.ResponseWriter, *http.Request)) {
	a.router.Post(path, http.HandlerFunc(f))
}

// Put wraps api router Put call
func (a *api) Put(path string, f func(http.ResponseWriter, *http.Request)) {
	a.router.Put(path, http.HandlerFunc(f))
}

// Delete wraps api router Delete call
func (a *api) Delete(path string, f func(http.ResponseWriter, *http.Request)) {
	a.router.Del(path, http.HandlerFunc(f))
}

// Options wraps api router Options call
func (a *api) Options(path string, f func(http.ResponseWriter, *http.Request)) {
	a.router.Options(path, http.HandlerFunc(f))
}

// Get wraps api router Get call
func (a *api) Get(path string, f func(http.ResponseWriter, *http.Request)) {
	a.router.Get(path, http.HandlerFunc(f))
}

// Head wraps api router Head call
func (a *api) Head(path string, f func(http.ResponseWriter, *http.Request)) {
	a.router.Head(path, http.HandlerFunc(f))
}

// AddHandler appends handler to api handlers
func (a *api) AddHandler(f func(http.ResponseWriter, *http.Request)) {
	a.handlers = append(a.handlers, f)
}

// Start initializes the api by setting up c3pio routes and robeaux
func (a *api) Start() {
	mcpCommandRoute := "/api/commands/:command"
	deviceCommandRoute := "/api/robots/:robot/devices/:device/commands/:command"
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
	a.Get(deviceCommandRoute, a.executeDeviceCommand)
	a.Post(deviceCommandRoute, a.executeDeviceCommand)
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
func (a *api) robeaux(res http.ResponseWriter, req *http.Request) {
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
	}
	res.Write(buf)
}

// mcp returns MCP route handler.
// Writes JSON with gobot representation
func (a *api) mcp(res http.ResponseWriter, req *http.Request) {
	a.writeJSON(map[string]interface{}{"MCP": a.gobot.ToJSON()}, res)
}

// mcpCommands returns commands route handler.
// Writes JSON with global commands representation
func (a *api) mcpCommands(res http.ResponseWriter, req *http.Request) {
	a.writeJSON(map[string]interface{}{"commands": a.gobot.ToJSON().Commands}, res)
}

// robots returns route handler.
// Writes JSON with robots representation
func (a *api) robots(res http.ResponseWriter, req *http.Request) {
	jsonRobots := []*gobot.JSONRobot{}
	a.gobot.Robots().Each(func(r *gobot.Robot) {
		jsonRobots = append(jsonRobots, r.ToJSON())
	})
	a.writeJSON(map[string]interface{}{"robots": jsonRobots}, res)
}

// robot returns route handler.
// Writes JSON with robot representation
func (a *api) robot(res http.ResponseWriter, req *http.Request) {
	a.writeJSON(map[string]interface{}{"robot": a.gobot.Robot(req.URL.Query().Get(":robot")).ToJSON()}, res)
}

// robotCommands returns commands route handler
// Writes JSON with robot commands representation
func (a *api) robotCommands(res http.ResponseWriter, req *http.Request) {
	a.writeJSON(map[string]interface{}{"commands": a.gobot.Robot(req.URL.Query().Get(":robot")).ToJSON().Commands}, res)
}

// robotDevices returns devices route handler.
// Writes JSON with robot devices representation
func (a *api) robotDevices(res http.ResponseWriter, req *http.Request) {
	jsonDevices := []*gobot.JSONDevice{}
	a.gobot.Robot(req.URL.Query().Get(":robot")).Devices().Each(func(d gobot.Device) {
		jsonDevices = append(jsonDevices, d.ToJSON())
	})
	a.writeJSON(map[string]interface{}{"devices": jsonDevices}, res)
}

// robotDevice returns device route handler.
// Writes JSON with robot device representation
func (a *api) robotDevice(res http.ResponseWriter, req *http.Request) {
	a.writeJSON(
		map[string]interface{}{"device": a.gobot.Robot(req.URL.Query().Get(":robot")).
			Device(req.URL.Query().Get(":device")).ToJSON()}, res,
	)
}

// robotDeviceEvent returns device event route handler.
// Creates an event stream connection
// and queries event data to be written when received
func (a *api) robotDeviceEvent(res http.ResponseWriter, req *http.Request) {
	f, _ := res.(http.Flusher)
	c, _ := res.(http.CloseNotifier)

	closer := c.CloseNotify()
	msg := make(chan string)

	res.Header().Set("Content-Type", "text/event-stream")
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Connection", "keep-alive")

	gobot.On(a.gobot.Robot(req.URL.Query().Get(":robot")).
		Device(req.URL.Query().Get(":device")).Event(req.URL.Query().Get(":event")),
		func(data interface{}) {
			d, _ := json.Marshal(data)
			msg <- string(d)
		})

	for {
		select {
		case data := <-msg:
			fmt.Fprintf(res, "data: %v\n\n", data)
			f.Flush()
		case <-closer:
			log.Println("Closing connection")
			return
		}
	}
}

// robotDeviceCommands returns device commands route handler
// writes JSON with robot device commands representation
func (a *api) robotDeviceCommands(res http.ResponseWriter, req *http.Request) {
	a.writeJSON(
		map[string]interface{}{"commands": a.gobot.Robot(req.URL.Query().Get(":robot")).
			Device(req.URL.Query().Get(":device")).ToJSON().Commands}, res,
	)
}

// robotConnections returns connections route handler
// writes JSON with robot connections representation
func (a *api) robotConnections(res http.ResponseWriter, req *http.Request) {
	jsonConnections := []*gobot.JSONConnection{}
	a.gobot.Robot(req.URL.Query().Get(":robot")).Connections().Each(func(c gobot.Connection) {
		jsonConnections = append(jsonConnections, c.ToJSON())
	})
	a.writeJSON(map[string]interface{}{"connections": jsonConnections}, res)
}

// robotConnection returns connection route handler
// writes JSON with robot connection representation
func (a *api) robotConnection(res http.ResponseWriter, req *http.Request) {
	a.writeJSON(
		map[string]interface{}{"connection": a.gobot.Robot(req.URL.Query().Get(":robot")).
			Connection(req.URL.Query().Get(":connection")).ToJSON()},
		res,
	)
}

// executeMcpCommand calls a global command asociated to requested route
func (a *api) executeMcpCommand(res http.ResponseWriter, req *http.Request) {
	a.executeCommand(a.gobot.Command(req.URL.Query().Get(":command")),
		res,
		req,
	)
}

// executeDeviceCommand calls a device command asociated to requested route
func (a *api) executeDeviceCommand(res http.ResponseWriter, req *http.Request) {
	a.executeCommand(
		a.gobot.Robot(req.URL.Query().Get(":robot")).
			Device(req.URL.Query().Get(":device")).
			Command(req.URL.Query().Get(":command")),
		res,
		req,
	)
}

// executeRobotCommand calls a robot command asociated to requested route
func (a *api) executeRobotCommand(res http.ResponseWriter, req *http.Request) {
	a.executeCommand(
		a.gobot.Robot(req.URL.Query().Get(":robot")).
			Command(req.URL.Query().Get(":command")),
		res,
		req,
	)
}

// executeCommand writes JSON response with `f` returned value.
func (a *api) executeCommand(f func(map[string]interface{}) interface{},
	res http.ResponseWriter,
	req *http.Request,
) {

	body := make(map[string]interface{})
	json.NewDecoder(req.Body).Decode(&body)

	if f != nil {
		a.writeJSON(map[string]interface{}{"result": f(body)}, res)
	} else {
		a.writeJSON("Unknown Command", res)
	}
}

// writeJSON writes `j` as JSON in response
func (a *api) writeJSON(j interface{}, res http.ResponseWriter) {
	data, _ := json.Marshal(j)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

// Debug add handler to api that prints each request
func (a *api) Debug() {
	a.AddHandler(func(res http.ResponseWriter, req *http.Request) {
		log.Println(req)
	})
}
