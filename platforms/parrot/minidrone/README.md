# Parrot Minidrone

The Parrot Minidrone is very inexpensive quadcopter that is controlled using Bluetooth LE.


## How to Install
```
go get -d -u gobot.io/x/gobot/... && go install gobot.io/x/gobot/platforms/ble
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
		drone.On(drone.Event("battery"), func(data interface{}) {
			fmt.Printf("battery: %d\n", data)
		})

		drone.On(drone.Event("status"), func(data interface{}) {
			fmt.Printf("status: %d\n", data)
		})

		drone.On(drone.Event("flying"), func(data interface{}) {
			fmt.Println("flying!")
			gobot.After(5*time.Second, func() {
				fmt.Println("landing...")
				drone.Land()
				drone.Land()
			})
		})

		drone.On(drone.Event("landed"), func(data interface{}) {
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

The Parrot Minidrone is a Bluetooth LE device.

You need to know the BLE ID of the Minidrone you want to connect to. The Gobot BLE client adaptor also lets you connect by friendly name, aka "RS_1234".

### OSX

To run any of the Gobot BLE code you must use the `GODEBUG=cgocheck=0` flag in order to get around some of the issues in the CGo-based implementation.

For example:

    GODEBUG=cgocheck=0 go run examples/minidrone.go RS_1234

OSX uses its own Bluetooth ID system which is different from the IDs used on Linux. The code calls thru the XPC interfaces provided by OSX, so as a result does not need to run under sudo.

### Ubuntu

On Linux the BLE code will need to run as a root user account. The easiest way to accomplish this is probably to use `go build` to build your program, and then to run the requesting executable using `sudo`.

For example:

    go build examples/minidrone.go
    sudo ./bb8 RS_1234

### Windows

Hopefully coming soon...
