package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
)

func Hello(params map[string]interface{}) string {
	name := params["name"].(string)
	return fmt.Sprintf("hi %v", name)
}

func main() {
	master := gobot.GobotMaster()

	server := api.NewApi(master)
	server.Username = "gort"
	server.Password = "klatuu"
	server.Start()

	hello := gobot.NewRobot("hello", nil, nil, nil)
	hello.Commands = map[string]interface{}{"Hello": Hello}

	master.Robots = append(master.Robots, hello)

	master.Start()
}
