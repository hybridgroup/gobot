//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/platforms/firmata"
)

func main() {
	board := firmata.NewAdaptor(os.Args[1])
	sensor := aio.NewGroveSoundSensorDriver(board, "3", aio.WithSensorCyclicRead(500*time.Millisecond))

	work := func() {
		_ = sensor.On(aio.Data, func(data interface{}) {
			fmt.Println("sensor", data)
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{board},
		[]gobot.Device{sensor},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
