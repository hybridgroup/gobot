# Sphero Ollie

The Sphero Ollie is a toy robot from Sphero that is controlled using Bluetooth LE. For more information, go to [http://www.sphero.com/ollie](http://www.sphero.com/ollie)

## How to Install

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use

```go
package main

import (
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/sphero/ollie"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	ollie := ollie.NewDriver(bleAdaptor)

	work := func() {
		gobot.Every(1*time.Second, func() {
			r := uint8(gobot.Rand(255))
			g := uint8(gobot.Rand(255))
			b := uint8(gobot.Rand(255))
			ollie.SetRGB(r, g, b)
		})
	}

	robot := gobot.NewRobot("ollieBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ollie},
		work,
	)

	robot.Start()
}
```

## How to Connect

The Sphero Ollie is a Bluetooth LE device.

You need to know the BLE ID of the Ollie you want to connect to. The Gobot BLE client adaptor also lets you connect by friendly name, aka "2B-1247".

### OSX

If you connect by name, then you do not need to worry about the Bluetooth LE ID. However, if you want to connect by ID, OS X uses its own Bluetooth ID system which is different from the IDs used on Linux. The code calls thru the XPC interfaces provided by OSX, so as a result does not need to run under sudo.

For example:

    go run examples/ollie.go 2B-1247

### Ubuntu

On Linux the BLE code will need to run as a root user account. The easiest way to accomplish this is probably to use `go build` to build your program, and then to run the requesting executable using `sudo`.

For example:

    go build examples/ollie.go
    sudo ./minidrone 2B-1247

### Windows

Hopefully coming soon...
