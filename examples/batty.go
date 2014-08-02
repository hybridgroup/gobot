package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
)

func main() {
	gbot := gobot.NewGobot()

	api.NewAPI(gbot).Start()

	gbot.AddCommand("echo", func(params map[string]interface{}) interface{} {
		return params["a"]
	})

	loopback := gobot.NewLoopbackAdaptor("loopback")
	ping := gobot.NewPingDriver(loopback, "ping")

	work := func() {
		gobot.Every(5*time.Second, func() {
			fmt.Println(ping.Ping())
		})
	}
	r := gobot.NewRobot("TestBot",
		[]gobot.Connection{loopback},
		[]gobot.Device{ping},
		work,
	)

	r.AddCommand("hello", func(params map[string]interface{}) interface{} {
		return fmt.Sprintf("Hello, %v!", params["greeting"])
	})

	gbot.AddRobot(r)
	gbot.Start()
}
