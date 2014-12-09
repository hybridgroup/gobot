package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/spark"
)

func main() {
	gbot := gobot.NewGobot()

	sparkCore := spark.NewSparkCoreAdaptor("spark", "DEVICE_ID", "ACCESS_TOKEN")

	work := func() {
		result, err := sparkCore.Function("brew", "hello")

		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("result from executing function is: %v", result)
		}
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{sparkCore},
		[]gobot.Device{},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
