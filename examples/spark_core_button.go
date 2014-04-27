package main

import (
	"fmt"
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

	button := gpio.NewButton(sparkCore)
	button.Name = "button"
	button.Pin = "D5"
	button.Interval = "2s"

	led := gpio.NewLed(sparkCore)
	led.Name = "led"
	led.Pin = "D7"

	work := func() {
		gobot.On(button.Events["push"], func(data interface{}) {
			led.On()
		})

		gobot.On(button.Events["release"], func(data interface{}) {
			led.Off()
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{sparkCore},
		Devices:     []gobot.Device{button, led},
		Work:        work,
	}

	robot.Start()
}
