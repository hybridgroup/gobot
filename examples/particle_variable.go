package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/particle"
)

func main() {
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

	robot.Start()
}
