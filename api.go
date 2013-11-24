package gobot

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"net/http"
)

func Api(bot *Gobot) {
	m := martini.Classic()

	m.Get("/robots", func() string {
		return toJson(bot.Robots)
	})

	m.Get("/robots/:robotname", func(params martini.Params) string {
		return toJson(bot.FindRobot(params["robotname"]))
	})

	m.Get("/robots/:robotname/devices", func(params martini.Params) string {
		return toJson(bot.FindRobot(params["robotname"]).GetDevices())
	})

	m.Get("/robots/:robotname/devices/:devicename", func(params martini.Params) string {
		return toJson(bot.FindRobotDevice(params["robotname"], params["devicename"]))
	})

	m.Post("/robots/:robotname/devices/:devicename/commands/:command", func(params martini.Params, res http.ResponseWriter, req *http.Request) string {
		decoder := json.NewDecoder(req.Body)
		var response_hash map[string]interface{}
		decoder.Decode(&response_hash)
		robot := bot.FindRobotDevice(params["robotname"], params["devicename"])
		p := make([]interface{}, len(response_hash))
		i := 0
		for _, v := range response_hash {
			p[i] = v
			i++
		}
		Call(robot.Driver, params["command"], p...)
		return ""
	})

	go m.Run()
}
