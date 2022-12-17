// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/platforms/tinkerboard"
)

// Wiring
// PWR  Tinkerboard: 1 (+3.3V, VCC), 2(+5V), 6, 9, 14, 20 (GND)
// SPI0 Tinkerboard (not working with armbian): 11 (CLK), 13 (TXD), 15 (RXD), 29 (CSN0), 31 (CSN1, n.c.)
// SPI2 Tinkerboard: 23 (CLK), 19 (TXD), 21 (RXD), 24 (CSN0), 26 (CSN1, n.c.)
// MFRC522 plate: VCC, GND, SCK (CLK), MOSI (->TXD), MISO (->RXD), NSS/SDA (CSN0/CSN1?)
func main() {
	a := tinkerboard.NewAdaptor()
	d := spi.NewMFRC522Driver(a, spi.WithBusNumber(2))

	work := func() {
		var err error

		gobot.Every(2*time.Second, func() {
			if err = d.Check(); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("spiBot",
		[]gobot.Connection{a},
		[]gobot.Device{d},
		work,
	)

	robot.Start()
}
