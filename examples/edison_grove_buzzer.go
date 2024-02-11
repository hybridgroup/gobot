//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/intel-iot/edison"
)

func main() {
	board := edison.NewAdaptor()
	buzzer := gpio.NewBuzzerDriver(board, "3")

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
		[]gobot.Connection{board},
		[]gobot.Device{buzzer},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
