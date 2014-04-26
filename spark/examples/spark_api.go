package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-gpio"
	"github.com/hybridgroup/gobot-spark"
)

func main() {
	master := gobot.GobotMaster()
	gobot.Api(master)

	spark := new(gobotSpark.SparkAdaptor)
	spark.Name = "spark"
	spark.Params = map[string]interface{}{
		"device_id":    "",
		"access_token": "",
	}

	led := gobotGPIO.NewLed(spark)
	led.Name = "led"
	led.Pin = "D7"

	work := func() {
		gobot.Every("1s", func() {
			led.Toggle()
		})
	}

	master.Robots = append(master.Robots, &gobot.Robot{
		Name:        "spark",
		Connections: []gobot.Connection{spark},
		Devices:     []gobot.Device{led},
		Work:        work,
	})

	master.Start()
}
