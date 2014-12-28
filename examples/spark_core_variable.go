package main

import (
	"fmt"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/spark"
)

func main() {
	gbot := gobot.NewGobot()

	sparkCore := spark.NewSparkCoreAdaptor("spark", "DEVICE_ID", "ACCESS_TOKEN")

	work := func() {
		gobot.Every(1*time.Second, func() {
			if temp, err := sparkCore.Variable("temperature"); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("result from \"temperature\" is:", temp)
			}
		})
	}

	robot := gobot.NewRobot("spark",
		[]gobot.Connection{sparkCore},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
