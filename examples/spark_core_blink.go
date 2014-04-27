package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gpio"
	"github.com/hybridgroup/gobot/spark"
)

func main() {
	sparkCore := spark.NewSparkCoreAdaptor()
	sparkCore.Name = "spark"
	sparkCore.Params = map[string]interface{}{
		"device_id":    "",
		"access_token": "",
	}

	led := gpio.NewLed(sparkCore)
	led.Name = "led"
	led.Pin = "D7"

	work := func() {
		gobot.Every("2s", func() {
			led.Toggle()
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{sparkCore},
		Devices:     []gobot.Device{led},
		Work:        work,
	}

	robot.Start()
}
