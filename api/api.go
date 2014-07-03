package api

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"github.com/bmizerany/pat"
	"github.com/hybridgroup/gobot"
	"io/ioutil"
	"log"
	"net/http"
)

// Optional restful API through Gobot has access
// all the robots.
type api struct {
	gobot    *gobot.Gobot
	server   *pat.PatternServeMux
	Host     string
	Port     string
	Username string
	Password string
	Cert     string
	Key      string
	Debug    bool
	start    func(*api)
}

func NewAPI(g *gobot.Gobot) *api {
	return &api{
		gobot: g,
		start: func(a *api) {
			if a == nil {
				return
			}

			port := a.Port
			if port == "" {
				port = "3000"
			}

			host := a.Host
			cert := a.Cert
			key := a.Key

			log.Println("Initializing API on " + host + ":" + port + "...")
			http.Handle("/", a.server)
			go func() {
				if cert != "" && key != "" {
					http.ListenAndServeTLS(host+":"+port, cert, key, nil)
				} else {
					log.Println("WARNING: API using insecure connection. We recommend using an SSL certificate with Gobot.")
					http.ListenAndServe(host+":"+port, nil)
				}
			}()
		},
	}
}

// start starts the api using the start function
// sets on the API on initialization.
func (a *api) Start() {
	a.server = pat.New()

	commandRoute := "/commands/:command"
	deviceCommandRoute := "/robots/:robot/devices/:device/commands/:command"
	robotCommandRoute := "/robots/:robot/commands/:command"

	a.server.Get("/", a.setHeaders(a.root))
	a.server.Get("/commands", a.setHeaders(a.commands))
	a.server.Get(commandRoute, a.setHeaders(a.executeCommand))
	a.server.Post(commandRoute, a.setHeaders(a.executeCommand))
	a.server.Get("/robots", a.setHeaders(a.robots))
	a.server.Get("/robots/:robot", a.setHeaders(a.robot))
	a.server.Get("/robots/:robot/commands", a.setHeaders(a.robotCommands))
	a.server.Get(robotCommandRoute, a.setHeaders(a.executeRobotCommand))
	a.server.Post(robotCommandRoute, a.setHeaders(a.executeRobotCommand))
	a.server.Get("/robots/:robot/devices", a.setHeaders(a.robotDevices))
	a.server.Get("/robots/:robot/devices/:device", a.setHeaders(a.robotDevice))
	a.server.Get("/robots/:robot/devices/:device/commands", a.setHeaders(a.robotDeviceCommands))
	a.server.Get(deviceCommandRoute, a.setHeaders(a.executeDeviceCommand))
	a.server.Post(deviceCommandRoute, a.setHeaders(a.executeDeviceCommand))
	a.server.Get("/robots/:robot/connections", a.setHeaders(a.robotConnections))
	a.server.Get("/robots/:robot/connections/:connection", a.setHeaders(a.robotConnection))

	a.start(a)
}

// basic auth inspired by https://github.com/codegangsta/martini-contrib/blob/master/auth/
func (a *api) basicAuth(res http.ResponseWriter, req *http.Request) bool {
	auth := req.Header.Get("Authorization")
	if !a.secureCompare(auth, "Basic "+base64.StdEncoding.EncodeToString([]byte(a.Username+":"+a.Password))) {
		res.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
		http.Error(res, "Not Authorized", http.StatusUnauthorized)
		return false
	}
	return true
}
func (a *api) secureCompare(given string, actual string) bool {
	if subtle.ConstantTimeEq(int32(len(given)), int32(len(actual))) == 1 {
		return subtle.ConstantTimeCompare([]byte(given), []byte(actual)) == 1
	}
	/* Securely compare actual to itself to keep constant time, but always return false */
	return subtle.ConstantTimeCompare([]byte(actual), []byte(actual)) == 1 && false
}

func (a *api) setHeaders(f func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if a.Debug {
			log.Println(req)
		}
		if a.Username != "" {
			if !a.basicAuth(res, req) {
				return
			}
		}
		f(res, req)
	}
}

func (a *api) root(res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(a.gobot.ToJSON())
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) commands(res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(a.gobot.ToJSON().Commands)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robots(res http.ResponseWriter, req *http.Request) {
	jsonRobots := []*gobot.JSONRobot{}
	a.gobot.Robots().Each(func(r *gobot.Robot) {
		jsonRobots = append(jsonRobots, r.ToJSON())
	})
	data, _ := json.Marshal(jsonRobots)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robot(res http.ResponseWriter, req *http.Request) {
	robot := req.URL.Query().Get(":robot")

	data, _ := json.Marshal(a.gobot.Robot(robot).ToJSON())
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robotCommands(res http.ResponseWriter, req *http.Request) {
	robot := req.URL.Query().Get(":robot")

	data, _ := json.Marshal(a.gobot.Robot(robot).ToJSON().Commands)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robotDevices(res http.ResponseWriter, req *http.Request) {
	robot := req.URL.Query().Get(":robot")

	jsonDevices := []*gobot.JSONDevice{}
	a.gobot.Robot(robot).Devices().Each(func(d gobot.Device) {
		jsonDevices = append(jsonDevices, d.ToJSON())
	})
	data, _ := json.Marshal(jsonDevices)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robotDevice(res http.ResponseWriter, req *http.Request) {
	robot := req.URL.Query().Get(":robot")
	device := req.URL.Query().Get(":device")

	data, _ := json.Marshal(a.gobot.Robot(robot).Device(device).ToJSON())
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robotDeviceCommands(res http.ResponseWriter, req *http.Request) {
	robot := req.URL.Query().Get(":robot")
	device := req.URL.Query().Get(":device")

	data, _ := json.Marshal(a.gobot.Robot(robot).Device(device).ToJSON().Commands)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robotConnections(res http.ResponseWriter, req *http.Request) {
	robot := req.URL.Query().Get(":robot")

	jsonConnections := []*gobot.JSONConnection{}
	a.gobot.Robot(robot).Connections().Each(func(c gobot.Connection) {
		jsonConnections = append(jsonConnections, c.ToJSON())
	})
	data, _ := json.Marshal(jsonConnections)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robotConnection(res http.ResponseWriter, req *http.Request) {
	robot := req.URL.Query().Get(":robot")
	connection := req.URL.Query().Get(":connection")

	data, _ := json.Marshal(a.gobot.Robot(robot).Connection(connection).ToJSON())
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) executeCommand(res http.ResponseWriter, req *http.Request) {
	command := req.URL.Query().Get(":command")

	data, _ := ioutil.ReadAll(req.Body)
	body := make(map[string]interface{})
	json.Unmarshal(data, &body)
	f := a.gobot.Commands()[command]

	if f != nil {
		data, _ = json.Marshal(f(body))
	} else {
		data, _ = json.Marshal("Unknown Command")
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) executeDeviceCommand(res http.ResponseWriter, req *http.Request) {
	robot := req.URL.Query().Get(":robot")
	device := req.URL.Query().Get(":device")
	command := req.URL.Query().Get(":command")

	data, _ := ioutil.ReadAll(req.Body)
	body := make(map[string]interface{})
	json.Unmarshal(data, &body)
	d := a.gobot.Robot(robot).Device(device)
	body["robot"] = robot
	f := d.Commands()[command]

	if f != nil {
		data, _ = json.Marshal(f(body))
	} else {
		data, _ = json.Marshal("Unknown Command")
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) executeRobotCommand(res http.ResponseWriter, req *http.Request) {
	robot := req.URL.Query().Get(":robot")
	command := req.URL.Query().Get(":command")

	data, _ := ioutil.ReadAll(req.Body)
	body := make(map[string]interface{})
	json.Unmarshal(data, &body)
	r := a.gobot.Robot(robot)
	body["robot"] = robot
	f := r.Commands[command]

	if f != nil {
		data, _ = json.Marshal(f(body))
	} else {
		data, _ = json.Marshal("Unknown Command")
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}
