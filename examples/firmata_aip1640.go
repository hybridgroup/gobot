//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_aip1640.go /dev/ttyACM0
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
	// In the WEMOS D1 Mini LED Matrix Shield clockPin = 14, dataPin = 13
	aip1640 := gpio.NewAIP1640Driver(firmataAdaptor, "14", "13")

	smiles := [3][8]byte{
		{0x3C, 0x42, 0xA5, 0x81, 0xA5, 0x99, 0x42, 0x3C}, // happy :)
		{0x3C, 0x42, 0xA5, 0x81, 0xBD, 0x81, 0x42, 0x3C}, // normal :|
		{0x3C, 0x42, 0xA5, 0x81, 0x99, 0xA5, 0x42, 0x3C}, // sad  :(
	}

	s := 0
	work := func() {
		aip1640.Clear()
		gobot.Every(600*time.Millisecond, func() {
			aip1640.DrawMatrix(smiles[s])
			if err := aip1640.Display(); err != nil {
				fmt.Println(err)
			}
			s++
			if s > 2 {
				s = 0
			}
		})
	}

	robot := gobot.NewRobot("aip1640Bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{aip1640},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
