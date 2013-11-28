package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
)

func Hello(params map[string]interface{}) string {
	name := params["name"].(string)
	hello.roll(90)
	return fmt.Sprintf("hi %v", name)
}

func main() {
	master := gobot.GobotMaster()
	gobot.Api(master)

	hello := new(gobot.Robot)
	hello.Name = "hello"
	hello.Work = func() {
	}
	hello.Commands = map[string]interface{}{"Hello": Hello}

	master.Robots = append(master.Robots, *hello)

	master.Start()
}
