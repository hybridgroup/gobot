package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/spark"
	"time"
)

func main() {
	master := gobot.NewGobot()
	api.NewApi(master).Start()

	sparkCore := spark.NewSparkCoreAdaptor("spark", "device_id", "access_token")
	led := gpio.NewLedDriver(sparkCore, "led", "D7")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	master.Robots = append(master.Robots,
		gobot.NewRobot("spark", []gobot.Connection{sparkCore}, []gobot.Device{led}, work))

	master.Start()
}
