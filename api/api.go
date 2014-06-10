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

func NewApi(g *gobot.Gobot) *api {
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
		a.robot_commands(params["robotname"], res, req)
	})

	robot_command_route := "/robots/:robotname/commands/:command"

	a.server.Get(robot_command_route, func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.executeRobotCommand(params["robotname"], params["command"], res, req)
	})
	a.server.Post(robot_command_route, func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.executeRobotCommand(params["robotname"], params["command"], res, req)
	})

	a.server.Get("/robots/:robotname/devices", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robot_devices(params["robotname"], res, req)
	})

	a.server.Get("/robots/:robotname/devices/:devicename", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robot_device(params["robotname"], params["devicename"], res, req)
	})

	a.server.Get("/robots/:robotname/devices/:devicename/commands", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robot_device_commands(params["robotname"], params["devicename"], res, req)
	})

	command_route := "/robots/:robotname/devices/:devicename/commands/:command"

	a.server.Get(command_route, func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.executeCommand(params["robotname"], params["devicename"], params["command"], res, req)
	})
	a.server.Post(command_route, func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.executeCommand(params["robotname"], params["devicename"], params["command"], res, req)
	})

	a.server.Get("/robots/:robotname/connections", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robot_connections(params["robotname"], res, req)
	})

	a.server.Get("/robots/:robotname/connections/:connectionname", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robot_connection(params["robotname"], params["connectionname"], res, req)
	})

	a.start(a)
}

func (a *api) robots(res http.ResponseWriter, req *http.Request) {
	jsonRobots := make([]*gobot.JsonRobot, 0)
	for _, robot := range a.gobot.Robots {
		jsonRobots = append(jsonRobots, robot.ToJson())
	}
	data, _ := json.Marshal(jsonRobots)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robot(name string, res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(a.gobot.Robot(name).ToJson())
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robot_commands(name string, res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(a.gobot.Robot(name).RobotCommands)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robot_devices(name string, res http.ResponseWriter, req *http.Request) {
	devices := a.gobot.Robot(name).Devices()
	jsonDevices := make([]*gobot.JsonDevice, 0)
	for _, device := range devices {
		jsonDevices = append(jsonDevices, device.ToJson())
	}
	data, _ := json.Marshal(jsonDevices)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robot_device(robot string, device string, res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(a.gobot.Robot(robot).Device(device).ToJson())
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robot_device_commands(robot string, device string, res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(a.gobot.Robot(robot).Device(device).Commands())
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robot_connections(name string, res http.ResponseWriter, req *http.Request) {
	connections := a.gobot.Robot(name).Connections()
	jsonConnections := make([]*gobot.JsonConnection, 0)
	for _, connection := range connections {
		jsonConnections = append(jsonConnections, connection.ToJson())
	}
	data, _ := json.Marshal(jsonConnections)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) robot_connection(robot string, connection string, res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(a.gobot.Robot(robot).Connection(connection).ToJson())
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
			ret := make([]interface{}, 0)
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
		ret := make([]interface{}, 0)
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
