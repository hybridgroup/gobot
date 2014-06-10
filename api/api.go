package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/go-martini/martini"
	"github.com/hybridgroup/gobot"
	"github.com/martini-contrib/auth"
	"github.com/martini-contrib/cors"
)

// Optional restful API through Gobot has access
// all the robots.
type api struct {
	gobot    *gobot.Gobot
	server   *martini.ClassicMartini
	Host     string
	Port     string
	Username string
	Password string
	Cert     string
	Key      string
	start    func(*api)
}

func NewAPI(g *gobot.Gobot) *api {
	return &api{
		gobot: g,
		start: func(a *api) {
			if a == nil {
				return
			}

			username := a.Username
			if username != "" {
				password := a.Password
				a.server.Use(auth.Basic(username, password))
			}

			port := a.Port
			if port == "" {
				port = "3000"
			}

			host := a.Host
			cert := a.Cert
			key := a.Key

			log.Println("Initializing API on " + host + ":" + port + "...")
			go func() {
				if cert != "" && key != "" {
					http.ListenAndServeTLS(host+":"+port, cert, key, a.server)
				} else {
					log.Println("WARNING: API using insecure connection. We recommend using an SSL certificate with Gobot.")
					http.ListenAndServe(host+":"+port, a.server)
				}
			}()
		},
	}
}

// start starts the api using the start function
// sets on the API on initialization.
func (a *api) Start() {
	a.server = martini.Classic()

	a.server.Use(martini.Static("robeaux"))
	a.server.Use(cors.Allow(&cors.Options{
		AllowAllOrigins: true,
	}))

	a.server.Get("/robots", func(res http.ResponseWriter, req *http.Request) {
		a.robots(res, req)
	})

	a.server.Get("/robots/:robotname", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robot(params["robotname"], res, req)
	})

	a.server.Get("/robots/:robotname/commands", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robotCommands(params["robotname"], res, req)
	})

	robotCommandRoute := "/robots/:robotname/commands/:command"

	a.server.Get(robotCommandRoute, func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.executeRobotCommand(params["robotname"], params["command"], res, req)
	})
	a.server.Post(robotCommandRoute, func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.executeRobotCommand(params["robotname"], params["command"], res, req)
	})

	a.server.Get("/robots/:robotname/devices", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robotDevices(params["robotname"], res, req)
	})

	a.server.Get("/robots/:robotname/devices/:devicename", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robotDevice(params["robotname"], params["devicename"], res, req)
	})

	a.server.Get("/robots/:robotname/devices/:devicename/commands", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robotDeviceCommands(params["robotname"], params["devicename"], res, req)
	})

	commandRoute := "/robots/:robotname/devices/:devicename/commands/:command"

	a.server.Get(commandRoute, func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.executeCommand(params["robotname"], params["devicename"], params["command"], res, req)
	})
	a.server.Post(commandRoute, func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.executeCommand(params["robotname"], params["devicename"], params["command"], res, req)
	})

	a.server.Get("/robots/:robotname/connections", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robotConnections(params["robotname"], res, req)
	})

	a.server.Get("/robots/:robotname/connections/:connectionname", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robotConnection(params["robotname"], params["connectionname"], res, req)
	})

	a.start(a)
}

func (a *api) robots(res http.ResponseWriter, req *http.Request) {
	jsonRobots := []*gobot.JSONRobot{}
	for _, robot := range a.gobot.Robots {
		jsonRobots = append(jsonRobots, robot.ToJSON())
	}
	data, _ := json.Marshal(jsonRobots)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robot(name string, res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(a.gobot.Robot(name).ToJSON())
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robotCommands(name string, res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(a.gobot.Robot(name).RobotCommands)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robotDevices(name string, res http.ResponseWriter, req *http.Request) {
	devices := a.gobot.Robot(name).Devices()
	jsonDevices := []*gobot.JSONDevice{}
	for _, device := range devices {
		jsonDevices = append(jsonDevices, device.ToJSON())
	}
	data, _ := json.Marshal(jsonDevices)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robotDevice(robot string, device string, res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(a.gobot.Robot(robot).Device(device).ToJSON())
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robotDeviceCommands(robot string, device string, res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(a.gobot.Robot(robot).Device(device).Commands())
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robotConnections(name string, res http.ResponseWriter, req *http.Request) {
	connections := a.gobot.Robot(name).Connections()
	jsonConnections := []*gobot.JSONConnection{}
	for _, connection := range connections {
		jsonConnections = append(jsonConnections, connection.ToJSON())
	}
	data, _ := json.Marshal(jsonConnections)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robotConnection(robot string, connection string, res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(a.gobot.Robot(robot).Connection(connection).ToJSON())
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) executeCommand(robotname string, devicename string, commandname string, res http.ResponseWriter, req *http.Request) {
	data, _ := ioutil.ReadAll(req.Body)
	var body map[string]interface{}
	json.Unmarshal(data, &body)
	robot := a.gobot.Robot(robotname).Device(devicename)
	commands := robot.Commands().([]string)
	for command := range commands {
		if commands[command] == commandname {
			ret := []interface{}{}
			for _, v := range gobot.Call(robot.Driver, commandname, body) {
				ret = append(ret, v.Interface())
			}
			data, _ = json.Marshal(ret)
			res.Header().Set("Content-Type", "application/json; charset=utf-8")
			res.Write(data)
			return
		}
	}
	data, _ = json.Marshal([]interface{}{"Unknown Command"})
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) executeRobotCommand(robotname string, commandname string, res http.ResponseWriter, req *http.Request) {
	data, _ := ioutil.ReadAll(req.Body)
	body := make(map[string]interface{})
	json.Unmarshal(data, &body)
	robot := a.gobot.Robot(robotname)
	in := make([]reflect.Value, 1)
	body["robotname"] = robotname
	in[0] = reflect.ValueOf(body)
	command := robot.Commands[commandname]
	if command != nil {
		ret := []interface{}{}
		for _, v := range reflect.ValueOf(robot.Commands[commandname]).Call(in) {
			ret = append(ret, v.Interface())
		}
		data, _ = json.Marshal(ret)
	} else {
		data, _ = json.Marshal([]interface{}{"Unknown Command"})
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}
