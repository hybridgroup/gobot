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
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/intel-iot/joule"
)

func main() {
	e := joule.NewAdaptor()
	ads1015 := i2c.NewADS1015Driver(e)
	sensor := aio.NewAnalogSensorDriver(ads1015, "0", aio.WithSensorCyclicRead(500*time.Millisecond))
	led := gpio.NewLedDriver(e, "J12_26")

	work := func() {
		_ = sensor.On(aio.Data, func(data interface{}) {
			brightness := uint8(gobot.ToScale(gobot.FromScale(float64(data.(int)), 0, 1023), 0, 255))
			fmt.Println("sensor", data)
			fmt.Println("brightness", brightness)
			if err := led.Brightness(brightness); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{e},
		[]gobot.Device{ads1015, sensor, led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
