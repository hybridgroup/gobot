package main

import (
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

	led := gobotGPIO.NewLed(spark)
	led.Name = "led"
	led.Pin = "A1"

	work := func() {
		brightness := uint8(0)
		fade_amount := uint8(15)

		gobot.Every("0.5s", func() {
			led.Brightness(brightness)
			brightness = brightness + fade_amount
			if brightness == 0 || brightness == 255 {
				fade_amount = -fade_amount
			}
		})
	}

	robot := gobot.Robot{
		Connections: []Connection{spark},
		Devices:     []Device{led},
		Work:        work,
	}

	robot.Start()
}
