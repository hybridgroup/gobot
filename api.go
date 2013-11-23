package gobot

import "github.com/codegangsta/martini"

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

	go m.Run()
}
