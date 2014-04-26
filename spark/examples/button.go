package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-gpio"
	"github.com/hybridgroup/gobot-spark"
)

func main() {

	spark := new(gobotSpark.SparkAdaptor)
	spark.Name = "spark"
	spark.Params = map[string]interface{}{
		"device_id":    "",
		"access_token": "",
	}

	button := gobotGPIO.NewButton(spark)
	button.Name = "button"
	button.Pin = "D5"
	button.Interval = "2s"

	led := gobotGPIO.NewLed(spark)
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
		Connections: []gobot.Connection{spark},
		Devices:     []gobot.Device{button, led},
		Work:        work,
	}

	robot.Start()
}
