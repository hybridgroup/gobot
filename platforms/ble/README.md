# Bluetooth LE

The Gobot BLE adaptor makes it easy to interact with Bluetooth LE aka Bluetooth 4.0 using Go.

It is written using the [ble](https://github.com/currantlabs/ble) package by [@roylee17](https://github.com/roylee17). Thank you!

Learn more about Bluetooth LE at http://en.wikipedia.org/wiki/Bluetooth_low_energy

This package also includes drivers for several well-known BLE Services:

- Battery Service
- Device Information Service
- Generic Access Service

## How to Install
```
go get -d -u gobot.io/x/gobot/...
```

### OSX

You need to have XCode installed to be able to compile code that uses the Gobot BLE adaptor on OSX. This is because the `ble` package uses a CGo based implementation.

### Ubuntu

Everything should already just compile on most Linux systems.

## How To Connect

When using BLE a "peripheral" aka "server" is something you connect to such a a pulse meter. A "central" aka "client" is what does the connecting, such as your computer or mobile phone.

You need to know the BLE ID of the peripheral you want to connect to. The Gobot BLE client adaptor also lets you connect to a peripheral by friendly name.

### OSX

If you connect by name, then you do not need to worry about the Bluetooth LE ID. However, if you want to connect by ID, OS X uses its own Bluetooth ID system which is different from the IDs used on Linux. The code calls thru the XPC interfaces provided by OSX, so as a result does not need to run under sudo.

For example:

    go run examples/minidrone.go 8b2f8032290143e18fc7c426619632e8

### Ubuntu

On Linux the BLE code will need to run as a root user account. The easiest way to accomplish this is probably to use `go build` to build your program, and then to run the requesting executable using `sudo`.

For example:

    go build examples/minidrone.go
    sudo ./minidrone AA:BB:CC:DD:EE

### Windows

Hopefully coming soon...

## How to Use

Here is an example that uses the BLE "Battery" service to retrieve the current change level of the peripheral device:

```go
package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	battery := ble.NewBatteryDriver(bleAdaptor)

	work := func() {
		gobot.Every(5*time.Second, func() {
			fmt.Println("Battery level:", battery.GetBatteryLevel())
		})
	}

	robot := gobot.NewRobot("bleBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{battery},
		work,
	)

	robot.Start()
}
```
