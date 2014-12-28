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
		if stream, err := sparkCore.EventStream("all", ""); err != nil {
			fmt.Println(err)
		} else {
			gobot.On(stream, func(data interface{}) {
				fmt.Println(data.(spark.Event))
			})
		}
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{sparkCore},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
