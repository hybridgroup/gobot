//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_buzzer.go /dev/ttyACM0
*/

package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])
	buzzer := gpio.NewBuzzerDriver(firmataAdaptor, "3")

	work := func() {
		type note struct {
			tone     float64
			duration float64
		}

		song := []note{
			{gpio.C4, gpio.Quarter},
			{gpio.C4, gpio.Quarter},
			{gpio.G4, gpio.Quarter},
			{gpio.G4, gpio.Quarter},
			{gpio.A4, gpio.Quarter},
			{gpio.A4, gpio.Quarter},
			{gpio.G4, gpio.Half},
			{gpio.F4, gpio.Quarter},
			{gpio.F4, gpio.Quarter},
			{gpio.E4, gpio.Quarter},
			{gpio.E4, gpio.Quarter},
			{gpio.D4, gpio.Quarter},
			{gpio.D4, gpio.Quarter},
			{gpio.C4, gpio.Half},
		}

		for _, val := range song {
			if err := buzzer.Tone(val.tone, val.duration); err != nil {
				fmt.Println(err)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{buzzer},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
