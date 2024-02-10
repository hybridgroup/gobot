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
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

func main() {
	board := raspi.NewAdaptor()
	ads1015 := i2c.NewADS1015Driver(board)
	sensor := aio.NewGroveRotaryDriver(ads1015, "0", aio.WithSensorCyclicRead(500*time.Millisecond))

	work := func() {
		_ = sensor.On(aio.Data, func(data interface{}) {
			fmt.Println("sensor", data)
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{board},
		[]gobot.Device{ads1015, sensor},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
