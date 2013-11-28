package gobot

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"net/http"
	"reflect"
)

type api struct{}

func Api(bot *Master) {
	a := new(api)
	m := martini.Classic()

	m.Get("/robots", func() string {
		return toJson(bot.Robots)
	})

	m.Get("/robots/:robotname", func(params martini.Params) string {
		return toJson(bot.FindRobot(params["robotname"]))
	})

	m.Get("/robots/:robotname/commands", func(params martini.Params) string {
		return toJson(bot.FindRobot(params["robotname"]).RobotCommands)
	})

	robot_command_route := "/robots/:robotname/commands/:command"

	m.Get(robot_command_route, func(params martini.Params, res http.ResponseWriter, req *http.Request) string {
		decoder := json.NewDecoder(req.Body)
		var body map[string]interface{}
		decoder.Decode(&body)
		return a.executeRobotCommand(bot, params, body)
	})

	m.Get("/robots/:robotname/devices", func(params martini.Params) string {
		return toJson(bot.FindRobot(params["robotname"]).GetDevices())
	})

	m.Get("/robots/:robotname/devices/:devicename", func(params martini.Params) string {
		return toJson(bot.FindRobotDevice(params["robotname"], params["devicename"]))
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
