package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/beaglebone"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
	gbot := gobot.NewGobot()

	beagleboneAdaptor := beaglebone.NewBeagleboneAdaptor("beaglebone")
	sensor := gpio.NewAnalogSensorDriver(beagleboneAdaptor, "sensor", "P9_33")
	led := gpio.NewLedDriver(beagleboneAdaptor, "led", "P9_14")

	work := func() {
		gobot.On(sensor.Event("data"), func(data interface{}) {
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
