package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/beaglebone"
)

func main() {
	a := beaglebone.NewAdaptor()
	sensor := gpio.NewAnalogSensorDriver(a, "P9_39")

	work := func() {
		sensor.On(gpio.Data, func(data interface{}) {
			voltage := (float64(data.(int)) * 1.8) / 1024 // BBB uses 1.8V
			tempC := (voltage - 0.5) * 100
			tempF := (tempC * 9 / 5) + 32

			fmt.Printf("%.2f°C\n", tempC)
			fmt.Printf("%.2f°F\n", tempF)
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{a},
		[]gobot.Device{sensor},
		work,
	)

	robot.Start()
}
