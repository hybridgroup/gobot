# GoPiGo3

The GoPiGo3 is a robotics controller by Dexter Industries that is compatible with the Raspberry Pi.

## How to Install

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use
This example will blink the left and right leds red/blue.

```go
package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	g "gobot.io/x/gobot/platforms/dexter/gopigo3"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	raspiAdaptor := raspi.NewAdaptor()
	gopigo3 := g.NewDriver(raspiAdaptor)

	work := func() {
		on := uint8(0xFF)
		gobot.Every(1000*time.Millisecond, func() {
			err := gopigo3.SetLED(g.LED_EYE_RIGHT, 0x00, 0x00, on)
			if err != nil {
				fmt.Println(err)
			}
			err = gopigo3.SetLED(g.LED_EYE_LEFT, ^on, 0x00, 0x00)
			if err != nil {
				fmt.Println(err)
			}
			on = ^on
		})
	}

	robot := gobot.NewRobot("gopigo3",
		[]gobot.Connection{raspiAdaptor},
		[]gobot.Device{gopigo3},
		work,
	)

	robot.Start()
}
```