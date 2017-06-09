# Parrot Minidrone

The Parrot Minidrones are very inexpensive drones that are controlled using Bluetooth LE aka Bluetooth 4.0.

Models that are known to work with this package include:

	- Parrot Rolling Spider
	- Parrot Airborne Cargo Mars
	- Parrot Airborne Cargo Travis
	- Parrot Mambo

Models that should work now, but have not been tested by us:

	- Parrot Airborne Night Swat
	- Parrot Airborne Night Maclane
	- Parrot Airborne Night Blaze
	- Parrot HYDROFOIL Orak
	- Parrot HYDROFOIL NewZ

Models that will require additional work for compatibility:

	- Parrot Swing

## How to Install

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use

```go
package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/parrot/minidrone"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	drone := minidrone.NewDriver(bleAdaptor)

	work := func() {
		drone.On(minidrone.Battery, func(data interface{}) {
			fmt.Printf("battery: %d\n", data)
		})

		drone.On(minidrone.FlightStatus, func(data interface{}) {
			fmt.Printf("flight status: %d\n", data)
		})

		drone.On(minidrone.Takeoff, func(data interface{}) {
			fmt.Println("taking off...")
		})

		drone.On(minidrone.Hovering, func(data interface{}) {
			fmt.Println("hovering!")
			gobot.After(5*time.Second, func() {
				drone.Land()
			})
		})

		drone.On(minidrone.Landing, func(data interface{}) {
			fmt.Println("landing...")
		})

		drone.On(minidrone.Landed, func(data interface{}) {
			fmt.Println("landed.")
		})

		time.Sleep(1000 * time.Millisecond)
		drone.TakeOff()
	}

	robot := gobot.NewRobot("minidrone",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}
```

## How to Connect

The Parrot Minidrones are Bluetooth LE devices.

You need to know the BLE ID or name of the Minidrone you want to connect to. The Gobot BLE client adaptor also lets you connect by friendly name, aka "RS_1234".

### OSX

If you connect by name, then you do not need to worry about the Bluetooth LE ID. However, if you want to connect by ID, OS X uses its own Bluetooth ID system which is different from the IDs used on Linux. The code calls thru the XPC interfaces provided by OSX, so as a result does not need to run under sudo.

For example:

    go run examples/minidrone.go RS_1234

### Ubuntu

On Linux the BLE code will need to run as a root user account. The easiest way to accomplish this is probably to use `go build` to build your program, and then to run the requesting executable using `sudo`.

For example:

    go build examples/minidrone.go
    sudo ./minidrone RS_1234

### Windows

Hopefully coming soon...
