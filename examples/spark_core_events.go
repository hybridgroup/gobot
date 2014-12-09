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
		stream, err := sparkCore.EventStream("all", "")

		if err != nil {
			fmt.Println(err.Error())
		} else {
			for {
				ev := <-stream.Events
				fmt.Println(ev.Event(), ev.Data())
			}
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
