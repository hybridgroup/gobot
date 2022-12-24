// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/platforms/adaptors"
	"gobot.io/x/gobot/platforms/tinkerboard"
)

// Wiring
// PWR  Tinkerboard: 1 (+3.3V, VCC), 2(+5V), 6, 9, 14, 20 (GND)
// GPIO-SPI Tinkerboard (same as SPI2): 23 (CLK), 19 (TXD), 21 (RXD), 24 (CSN0)
// MFRC522 plate: VCC, GND, SCK (CLK), MOSI (->TXD), MISO (->RXD), NSS/SDA (CSN0/CSN1?)
const (
	sclk    = "23"
	nss     = "24"
	mosi    = "19"
	miso    = "21"
	speedHz = 5000 // more than 15kHz is not possible with GPIO's, so we choose 5kHz
)

func main() {
	a := tinkerboard.NewAdaptor(adaptors.WithSpiGpioAccess(sclk, nss, mosi, miso))
	d := spi.NewMFRC522Driver(a, spi.WithSpeed(speedHz))

	wasCardDetected := false
	const textToCard = "Hello RFID user!\nHey, we use GPIO's for SPI."

	work := func() {
		if err := d.PrintReaderVersion(); err != nil {
			fmt.Println("get version err:", err)
		}

		gobot.Every(2*time.Second, func() {

			if !wasCardDetected {
				fmt.Println("\n+++ poll for card +++")
				if err := d.IsCardPresent(); err != nil {
					fmt.Println("no card found")
				} else {
					fmt.Println("\n+++ write card +++")
					err := d.WriteText(textToCard)
					if err != nil {
						fmt.Println("write err:", err)
					}
					wasCardDetected = true
				}
			} else {
				fmt.Println("\n+++ read card +++")
				text, err := d.ReadText()
				if err != nil {
					fmt.Println("read err:", err)
					wasCardDetected = false
				} else {
					fmt.Printf("-- start text --\n%s\n-- end  text --\n", text)
				}
			}

		})
	}

	robot := gobot.NewRobot("gpioSpiBot",
		[]gobot.Connection{a},
		[]gobot.Device{d},
		work,
	)

	robot.Start()
}
