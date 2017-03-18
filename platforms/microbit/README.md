# Microbit

The [Microbit](http://microbit.org/) is a tiny computer with built-in Bluetooth LE aka Bluetooth 4.0.

## How to Install
```
go get -d -u gobot.io/x/gobot/... && go install gobot.io/x/gobot/platforms/microbit
```

You must install the Microbit firmware from [@sandeepmistry] located at  [https://github.com/sandeepmistry/node-bbc-microbit](https://github.com/sandeepmistry/node-bbc-microbit) to use the Microbit with Gobot. This firmware is based on the micro:bit template, but with a few changes.

If you have the [Gort](https://gort.io) command line tool installed, you can use the following commands:

```
gort microbit download
gort microbit install /media/mysystem/MICROBIT
```

You can also follow the firmware installation instructions at [https://github.com/sandeepmistry/node-bbc-microbit#flashing-microbit-firmware](https://github.com/sandeepmistry/node-bbc-microbit#flashing-microbit-firmware).

The source code for the firmware is located at [https://github.com/sandeepmistry/node-bbc-microbit-firmware](https://github.com/sandeepmistry/node-bbc-microbit-firmware) however you do not need this source code to install the firmware using the installation instructions.

## How to Use

The Gobot platform for the Microbit includes several different drivers, each one corresponding to a different capability:

- AccelerometerDriver
- ButtonDriver
- LEDDriver
- MagnetometerDriver
- TemperatureDriver

The following example uses the LEDDriver:

```go
package main

import (
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/microbit"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	ubit := microbit.NewLEDDriver(bleAdaptor)

	work := func() {
		ubit.Blank()
		gobot.After(1*time.Second, func() {
			ubit.WriteText("Hello")
		})
		gobot.After(7*time.Second, func() {
			ubit.Smile()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ubit},
		work,
	)

	robot.Start()
}
```

## How to Connect

The Microbit is a Bluetooth LE device.

You need to know the BLE ID of the Microbit that you want to connect to.

### OSX

To run any of the Gobot BLE code you must use the `GODEBUG=cgocheck=0` flag in order to get around some of the issues in the CGo-based implementation.

For example:

    GODEBUG=cgocheck=0 go run examples/microbit_blink.go "BBC micro:bit"

OSX uses its own Bluetooth ID system which is different from the IDs used on Linux. The code calls thru the XPC interfaces provided by OSX, so as a result does not need to run under sudo.

### Ubuntu

On Linux the BLE code will need to run as a root user account. The easiest way to accomplish this is probably to use `go build` to build your program, and then to run the requesting executable using `sudo`.

For example:

    go build examples/microbit_blink.go
    sudo ./microbit_blink "BBC micro:bit"

### Windows

Hopefully coming soon...
