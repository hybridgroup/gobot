# Leap

The Leap Motion is a user-interface device that tracks the user's hand motions, and translates them into events that can control robots and physical computing hardware.

For more info about the Leap Motion platform click [Leap Motion](https://www.leapmotion.com/)

## How to Install

First install the [Leap Motion Software](https://www.leapmotion.com/setup)

Now you can install the package with:

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use

```go
package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/leap"
)

func main() {
	leapMotionAdaptor := leap.NewAdaptor("127.0.0.1:6437")
	l := leap.NewDriver(leapMotionAdaptor)

	work := func() {
		l.On(l.Event("message"), func(data interface{}) {
			fmt.Println(data.(leap.Frame))
		})
	}

	robot := gobot.NewRobot("leapBot",
		[]gobot.Connection{leapMotionAdaptor},
		[]gobot.Device{l},
		work,
	)

	robot.Start()
}
```

## How To Connect

### OSX

This driver works out of the box with the vanilla installation of the Leap Motion Software that you get in their [Setup Guide](https://www.leapmotion.com/setup).

The main steps are:

*   Run Leap Motion.app to open a websocket connection in port 6437.
*   Connect your Computer and Leap Motion Controller.
*   Connect to the device via Gobot.

### Ubuntu

The Linux download of the Leap Motion software can be obtained from [Leap Motion Dev Center](https://developer.leapmotion.com/downloads) (requires free signup)

The main steps are:

*   Run the leapd daemon, to open a websocket connection in port 6437.
*   Connect your computer and the Leap Motion controller
*   Connect to the device via Gobot
