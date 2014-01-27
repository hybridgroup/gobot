package gobot

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"net/http"
	"reflect"
)

type api struct{}

type jsonRobot struct {
	Name        string            `json:"name"`
	Commands    []string          `json:"commands"`
	Connections []*jsonConnection `json:"connections"`
	Devices     []*jsonDevice     `json:"devices"`
}

type jsonDevice struct {
	Name       string          `json:"name"`
	Driver     string          `json:"driver"`
	Connection *jsonConnection `json:"connection"`
	Commands   []string        `json:"commands"`
}

type jsonConnection struct {
	Name    string `json:"name"`
	Port    string `json:"port"`
	Adaptor string `json:"adaptor"`
}

func Api(bot *Master) {
	a := new(api)
	m := martini.Classic()

	m.Get("/robots", func() string {
		jsonRobots := make([]*jsonRobot, 0)
		for _, robot := range bot.Robots {
			jsonRobots = append(jsonRobots, a.formatJsonRobot(&robot))
		}
		return toJson(jsonRobots)
	})

	m.Get("/robots/:robotname", func(params martini.Params) string {
		return toJson(a.formatJsonRobot(bot.FindRobot(params["robotname"])))
	})

	m.Get("/robots/:robotname/commands", func(params martini.Params) string {
		return toJson(bot.FindRobot(params["robotname"]).RobotCommands)
	})

	robot_command_route := "/robots/:robotname/commands/:command"

	m.Get(robot_command_route, func(params martini.Params, res http.ResponseWriter, req *http.Request) string {
		decoder := json.NewDecoder(req.Body)
		var body map[string]interface{}
		decoder.Decode(&body)
		if len(body) == 0 {
			body = map[string]interface{}{}
		}
		body["robotname"] = params["robotname"]
		return a.executeRobotCommand(bot, params, body)
	})

	m.Get("/robots/:robotname/devices", func(params martini.Params) string {
		devices := bot.FindRobot(params["robotname"]).GetDevices()
		jsonDevices := make([]*jsonDevice, 0)
		for _, device := range devices {
			jsonDevices = append(jsonDevices, a.formatJsonDevice(device))
		}
		return toJson(jsonDevices)
	})

	m.Get("/robots/:robotname/devices/:devicename", func(params martini.Params) string {
		return toJson(a.formatJsonDevice(bot.FindRobotDevice(params["robotname"], params["devicename"])))
	})

	m.Get("/robots/:robotname/devices/:devicename/commands", func(params martini.Params) string {
		return toJson(bot.FindRobotDevice(params["robotname"], params["devicename"]).Commands())
	})

	command_route := "/robots/:robotname/devices/:devicename/commands/:command"

	m.Get(command_route, func(params martini.Params, res http.ResponseWriter, req *http.Request) string {
		return a.executeCommand(bot, params, res, req)
	})
	m.Post(command_route, func(params martini.Params, res http.ResponseWriter, req *http.Request) string {
		return a.executeCommand(bot, params, res, req)
	})

	go m.Run()
}

func (a *api) formatJsonRobot(robot *Robot) *jsonRobot {
	jsonRobot := new(jsonRobot)
	jsonRobot.Name = robot.Name
	jsonRobot.Commands = robot.RobotCommands
	jsonRobot.Connections = make([]*jsonConnection, 0)
	for _, device := range robot.devices {
		jsonDevice := a.formatJsonDevice(device)
		jsonRobot.Connections = append(jsonRobot.Connections, jsonDevice.Connection)
		jsonRobot.Devices = append(jsonRobot.Devices, jsonDevice)
	}
	return jsonRobot
}

func (a *api) formatJsonDevice(device *device) *jsonDevice {
	jsonDevice := new(jsonDevice)
	jsonDevice.Name = device.Name
	jsonDevice.Driver = FieldByNamePtr(device.Driver, "Name").Interface().(string)
	jsonDevice.Connection = new(jsonConnection)
	jsonDevice.Connection.Name = FieldByNamePtr(FieldByNamePtr(device.Driver, "Adaptor").Interface().(AdaptorInterface), "Name").Interface().(string)
	jsonDevice.Connection.Port = FieldByNamePtr(FieldByNamePtr(device.Driver, "Adaptor").Interface().(AdaptorInterface), "Port").Interface().(string)
	jsonDevice.Connection.Adaptor = FieldByNamePtr(device.Driver, "Adaptor").Type().Name()
	jsonDevice.Commands = FieldByNamePtr(device.Driver, "Commands").Interface().([]string)
	return jsonDevice
}

func (a *api) executeCommand(bot *Master, params martini.Params, res http.ResponseWriter, req *http.Request) string {
	decoder := json.NewDecoder(req.Body)
	var body map[string]interface{}
	decoder.Decode(&body)
	robot := bot.FindRobotDevice(params["robotname"], params["devicename"])
	commands := robot.Commands().([]string)
	for command := range commands {
		if commands[command] == params["command"] {
			ret := Call(robot.Driver, params["command"], body)
			return toJson(map[string]interface{}{"results": ret})
		}
	}
	return toJson(map[string]interface{}{"results": "Unknown Command"})
}

func (a *api) executeRobotCommand(bot *Master, m_params martini.Params, params ...interface{}) string {
	robot := bot.FindRobot(m_params["robotname"])
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	ret := reflect.ValueOf(robot.Commands[m_params["command"]]).Call(in)
	return toJson(map[string]interface{}{"results": ret[0].Interface()})
}
