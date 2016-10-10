package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/drivers/gpio"
	"github.com/hybridgroup/gobot/platforms/beaglebone"
)

func main() {
	gbot := gobot.NewGobot()

	beagleboneAdaptor := beaglebone.NewAdaptor()
	sensor := gpio.NewAnalogSensorDriver(beagleboneAdaptor, "P9_33")
	led := gpio.NewLedDriver(beagleboneAdaptor, "P9_14")

	work := func() {
		sensor.On(sensor.Event("data"), func(data interface{}) {
			brightness := uint8(
				gobot.ToScale(gobot.FromScale(float64(data.(int)), 0, 1024), 0, 255),
			)
			fmt.Println("sensor", data)
			fmt.Println("brightness", brightness)
			led.Brightness(brightness)
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{beagleboneAdaptor},
		[]gobot.Device{sensor, led},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
