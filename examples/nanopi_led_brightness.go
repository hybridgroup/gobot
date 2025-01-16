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
	"gobot.io/x/gobot/v2/platforms/friendlyelec/nanopi"
)

// Wiring
// PWR  NanoPi: 1, 17 (+3.3V, VCC); 2, 4 (+5V, VDD); 6, 9, 14, 20 (GND)
// GPIO NanoPi: the fourth header pin at inner USB side, count from USB side, is the PWM output
// LED: the PWM output is NOT able to drive a 20mA LED with full brightness (custom driver or low current LED is needed)
// Expected behavior: the LED fades in and out
func main() {
	r := nanopi.NewNeoAdaptor()
	led := gpio.NewLedDriver(r, "PWM")

	work := func() {
		brightness := uint8(0)
		fadeAmount := uint8(15)

		gobot.Every(100*time.Millisecond, func() {
			if err := led.Brightness(brightness); err != nil {
				fmt.Println(err)
			}
			brightness = brightness + fadeAmount
			if brightness == 0 || brightness == 255 {
				fadeAmount = -fadeAmount
			}
		})
	}

	robot := gobot.NewRobot("pwmBot",
		[]gobot.Connection{r},
		[]gobot.Device{led},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
