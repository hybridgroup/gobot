//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/beagleboard/beaglebone"
)

func main() {
	beagleboneAdaptor := beaglebone.NewAdaptor()
	sensor := aio.NewAnalogSensorDriver(beagleboneAdaptor, "P9_33", aio.WithSensorCyclicRead(500*time.Millisecond))
	led := gpio.NewLedDriver(beagleboneAdaptor, "P9_14")

	work := func() {
		_ = sensor.On(sensor.Event("data"), func(data interface{}) {
			brightness := uint8(
				gobot.ToScale(gobot.FromScale(float64(data.(int)), 0, 1024), 0, 255),
			)
			fmt.Println("sensor", data)
			fmt.Println("brightness", brightness)
			if err := led.Brightness(brightness); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{beagleboneAdaptor},
		[]gobot.Device{sensor, led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
