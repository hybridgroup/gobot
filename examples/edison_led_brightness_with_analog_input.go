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
	"gobot.io/x/gobot/v2/platforms/intel-iot/edison"
)

func main() {
	e := edison.NewAdaptor()
	sensor := aio.NewAnalogSensorDriver(e, "0", aio.WithSensorCyclicRead(500*time.Millisecond))
	led := gpio.NewLedDriver(e, "3")

	work := func() {
		_ = sensor.On(aio.Data, func(data interface{}) {
			brightness := uint8(
				gobot.ToScale(gobot.FromScale(float64(data.(int)), 0, 4096), 0, 255),
			)
			fmt.Println("sensor", data)
			fmt.Println("brightness", brightness)
			if err := led.Brightness(brightness); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{e},
		[]gobot.Device{sensor, led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
