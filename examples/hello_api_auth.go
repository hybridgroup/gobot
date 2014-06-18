package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
)

func main() {
	master := gobot.NewGobot()

	server := api.NewAPI(master)
	server.Username = "gort"
	server.Password = "klatuu"
	server.Start()

	hello := gobot.NewRobot("hello", nil, nil, nil)

	hello.AddCommand("HiThere", func(params map[string]interface{}) interface{} {
		return []string{fmt.Sprintf("Hey"), fmt.Sprintf("dude!")}
	})

	master.Robots = append(master.Robots, hello)

	master.Start()
}
