package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/particle"
)

func main() {
	gbot := gobot.NewMaster()

	core := particle.NewAdaptor("DEVICE_ID", "ACCESS_TOKEN")

	work := func() {
		gobot.Every(1*time.Second, func() {
			if temp, err := core.Variable("temperature"); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("result from \"temperature\" is:", temp)
			}
		})
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{core},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
