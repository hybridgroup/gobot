package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
)

func main() {
	gbot := gobot.NewGobot()

	server := api.NewAPI(gbot)
	server.Username = "gort"
	server.Password = "klatuu"
	server.Start()

	hello := gbot.AddRobot(gobot.NewRobot("hello"))

	hello.AddCommand("hi_there", func(params map[string]interface{}) interface{} {
		return []string{fmt.Sprintf("Hey"), fmt.Sprintf("dude!")}
	})

	gbot.Start()
}
