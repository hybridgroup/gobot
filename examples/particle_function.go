package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/particle"
)

func main() {
	core := particle.NewAdaptor("DEVICE_ID", "ACCESS_TOKEN")

	work := func() {
		if result, err := core.Function("brew", "202,230"); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("result from \"brew\":", result)
		}
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{core},
		work,
	)

	robot.Start()
}
