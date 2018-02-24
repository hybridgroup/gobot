// +build example
//
// Do not build by default.

/*
 How to setup
 This examples requires you to daisy-chain 4 led matrices based on either MAX7219 or MAX7221.
 It will turn on one led at a time, from the first led at the first matrix to the last led of the last matrix.

 How to run
 Pass serial port to use as the first param:

	go run examples/firmata_max72xx.go /dev/ttyACM0
*/

package main

import (
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor(os.Args[1])
	max := gpio.NewMAX72xxDriver(firmataAdaptor, "11", "10", "9", 4)

	var digit byte = 1 // digit address goes from 0x01 (MAX72xxDigit0) to 0x08 (MAX72xxDigit8)
	var bits byte = 1
	var module uint
	count := 0

	work := func() {
		gobot.Every(100*time.Millisecond, func() {
			max.ClearAll()
			max.One(module, digit, bits)
			bits = bits << 1

			count++
			if count > 7 {
				count = 0
				digit++
				bits = 1
				if digit > 8 {
					digit = 1
					module++
					if module >= 4 {
						module = 0
						count = 0
					}
				}
			}
		})
	}

	robot := gobot.NewRobot("Max72xxBot",
		[]gobot.Connection{esp8266},
		[]gobot.Device{max},
		work,
	)

	robot.Start()
}
