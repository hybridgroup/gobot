# Tello

This package contains the Gobot driver for the Ryze Tello drone, sold by DJI.

For more information on this drone, go to: https://www.ryzerobotics.com/tello

## How to Install

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use

Connect to the drone's Wi-Fi network from your computer. It will be named something like "TELLO-XXXXXX".

Once you are connected you can run the Gobot code on your computer to control the drone.

Here is a sample of how you initialize and use the driver:

```go
package main

import (
	"fmt"
    "time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

func main() {
	drone := tello.NewDriver("8888")

	work := func() {
		drone.TakeOff()

		gobot.After(5*time.Second, func() {
			drone.Land()
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}
```

## References

Thanks to https://github.com/microlinux/tello for serving as an example for the Tello community with his Python library.
