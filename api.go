package gobot

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
	"reflect"
)

type api struct {
	master *Master
	server *martini.ClassicMartini
	Host   string
	Port   string
}

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

var startApi = func(me *api) {
	port := me.Port
	if port == "" {
		port = "3000"
	}

	host := me.Host
	go http.ListenAndServe(host+":"+port, me.server)
}

func (me *api) StartApi() {
	startApi(me)
}

func Api(bot *Master) *api {
	a := new(api)
	a.master = bot
	bot.Api = a

	m := martini.Classic()
	a.server = m

	m.Use(martini.Static("robeaux"))

	m.Get("/robots", func(res http.ResponseWriter, req *http.Request) {
		a.robots(res, req)
	})

	m.Get("/robots/:robotname", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robot(params["robotname"], res, req)
	})

	m.Get("/robots/:robotname/commands", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robot_commands(params["robotname"], res, req)
	})

	robot_command_route := "/robots/:robotname/commands/:command"

	m.Get(robot_command_route, func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.executeRobotCommand(params["robotname"], params["command"], res, req)
	})
	m.Post(robot_command_route, func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.executeRobotCommand(params["robotname"], params["command"], res, req)
	})

	m.Get("/robots/:robotname/devices", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robot_devices(params["robotname"], res, req)
	})

	m.Get("/robots/:robotname/devices/:devicename", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robot_device(params["robotname"], params["devicename"], res, req)
	})

	m.Get("/robots/:robotname/devices/:devicename/commands", func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.robot_device_commands(params["robotname"], params["devicename"], res, req)
	})

	command_route := "/robots/:robotname/devices/:devicename/commands/:command"

	m.Get(command_route, func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.executeCommand(params["robotname"], params["devicename"], params["command"], res, req)
	})
	m.Post(command_route, func(params martini.Params, res http.ResponseWriter, req *http.Request) {
		a.executeCommand(params["robotname"], params["devicename"], params["command"], res, req)
	})

	return a
}

func (me *api) robots(res http.ResponseWriter, req *http.Request) {
	jsonRobots := make([]*jsonRobot, 0)
	for _, robot := range me.master.Robots {
		jsonRobots = append(jsonRobots, me.formatJsonRobot(robot))
	}
	data, _ := json.Marshal(jsonRobots)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (me *api) robot(name string, res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(me.formatJsonRobot(me.master.FindRobot(name)))
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (me *api) robot_commands(name string, res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(me.master.FindRobot(name).RobotCommands)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (me *api) robot_devices(name string, res http.ResponseWriter, req *http.Request) {
	devices := me.master.FindRobot(name).GetDevices()
	jsonDevices := make([]*jsonDevice, 0)
	for _, device := range devices {
		jsonDevices = append(jsonDevices, me.formatJsonDevice(device))
	}
	data, _ := json.Marshal(jsonDevices)
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (me *api) robot_device(robot string, device string, res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(me.formatJsonDevice(me.master.FindRobotDevice(robot, device)))
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (me *api) robot_device_commands(robot string, device string, res http.ResponseWriter, req *http.Request) {
	data, _ := json.Marshal(me.master.FindRobotDevice(robot, device).Commands())
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
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

func (a *api) executeCommand(robotname string, devicename string, commandname string, res http.ResponseWriter, req *http.Request) {
	data, _ := ioutil.ReadAll(req.Body)
	var body map[string]interface{}
	json.Unmarshal(data, &body)
	robot := a.master.FindRobotDevice(robotname, devicename)
	commands := robot.Commands().([]string)
	for command := range commands {
		if commands[command] == commandname {
			ret := Call(robot.Driver, commandname, body)
			data, _ = json.Marshal(map[string]interface{}{"result": ret[0].Interface()})
			res.Header().Set("Content-Type", "application/json; charset=utf-8")
			res.Write(data)
			return
		}
	}
	data, _ = json.Marshal(map[string]interface{}{"result": "Unknown Command"})
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}

func (a *api) executeRobotCommand(robotname string, commandname string, res http.ResponseWriter, req *http.Request) {
	data, _ := ioutil.ReadAll(req.Body)
	var body map[string]interface{}
	json.Unmarshal(data, &body)
	robot := a.master.FindRobot(robotname)
	in := make([]reflect.Value, 1)
	body["robotname"] = robotname
	in[0] = reflect.ValueOf(body)
	command := robot.Commands[commandname]
	if command != nil {
		ret := reflect.ValueOf(robot.Commands[commandname]).Call(in)
		data, _ = json.Marshal(map[string]interface{}{"result": ret[0].Interface()})
	} else {
		data, _ = json.Marshal(map[string]interface{}{"result": "Unknown Command"})
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Write(data)
}
