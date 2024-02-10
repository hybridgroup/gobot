//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_integration.go /dev/ttyACM0
*/

package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])
	led1 := gpio.NewLedDriver(firmataAdaptor, "3")
	led2 := gpio.NewLedDriver(firmataAdaptor, "4")
	button := gpio.NewButtonDriver(firmataAdaptor, "2")
	sensor := aio.NewAnalogSensorDriver(firmataAdaptor, "0", aio.WithSensorCyclicRead(500*time.Millisecond))

	work := func() {
		gobot.Every(1*time.Second, func() {
			if err := led1.Toggle(); err != nil {
				fmt.Println(err)
			}
		})
		gobot.Every(2*time.Second, func() {
			if err := led2.Toggle(); err != nil {
				fmt.Println(err)
			}
		})
		_ = button.On(gpio.ButtonPush, func(data interface{}) {
			if err := led2.On(); err != nil {
				fmt.Println(err)
			}
		})
		_ = button.On(gpio.ButtonRelease, func(data interface{}) {
			if err := led2.Off(); err != nil {
				fmt.Println(err)
			}
		})
		_ = sensor.On(aio.Data, func(data interface{}) {
			fmt.Println("sensor", data)
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led1, led2, button, sensor},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
