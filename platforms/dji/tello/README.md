# Tello

This package contains the Gobot driver for the Ryze Tello drone, sold by DJI.

For more information on this drone, go to: [https://www.ryzerobotics.com/tello](https://www.ryzerobotics.com/tello)

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

## Telo Edu driver

If you are planning to connect to the edu version of the tello, please use the `NewDriverWithIP` driver instead

```go
drone := tello.NewDriverWithIP("192.168.10.1", "8888")
```

## References

This driver could not exist without the awesome members of the unofficial Tello forum:

https://tellopilots.com/forums/tello-development.8/

Special thanks to [@Kragrathea](https://github.com/Kragrathea) who figured out a LOT of the packets and code as implemented in C#: [https://github.com/Kragrathea/TelloPC](https://github.com/Kragrathea/TelloPC)

Also thanks to [@microlinux](https://github.com/microlinux) with the Python library which served as the first example for the Tello community: [https://github.com/microlinux/tello](https://github.com/microlinux/tello)

Thank you to bluejune for the [https://bitbucket.org/PingguSoft/pytello](https://bitbucket.org/PingguSoft/pytello) repo, especially the Wireshark Lua dissector which has proven indispensable.
