// +build example
//
// Do not build by default.
/*
 * blinktest - blink the USR1 LED on the Beaglebone once a second.
 *
 * usr1 is the blue LED third from the Ethernet port in the group
 * of four.  Normally it goes on when the SD is being accessed.
 * While this is running expect it to blink on and off with a
 * two-second cycle. Afterward, the pin will be left unbound.
 *
 * Note: you must be running as root or be in the gpio group for
 * this to work.
 */
package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/beaglebone"

	"log"
)

func main() {
	log.Printf("Blink demo\n")
	pin := "usr1"
	beagleboneAdaptor := beaglebone.NewAdaptor()
	led := gpio.NewLedDriver(beagleboneAdaptor, pin)

	work := func() {
		gobot.Every(1*time.Second, func() {
			log.Printf("Toggling %s", pin)
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{beagleboneAdaptor},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}
