package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
)

func main() {
	gbot := gobot.NewGobot()
	api.NewAPI(gbot).Start()

	gbot.AddCommand("custom_gobot_command", func(params map[string]interface{}) interface{} {
		return "This command is attached to the master!"
	})

	hello := gbot.AddRobot(gobot.NewRobot("hello"))

	hello.AddCommand("hi_there", func(params map[string]interface{}) interface{} {
		return fmt.Sprintf("This command is attached to the robot %v", hello.Name)
	})

	gbot.Start()
}
