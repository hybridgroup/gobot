package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/particle"
)

func main() {
	core := particle.NewAdaptor("DEVICE_ID", "ACCESS_TOKEN")

	work := func() {
		if stream, err := core.EventStream("all", ""); err != nil {
			fmt.Println(err)
		} else {
			// TODO: some other way to handle this
			fmt.Println(stream)
		}
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{core},
		work,
	)

	robot.Start()
}
