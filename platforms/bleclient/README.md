# Bluetooth LE

The Gobot BLE adaptor makes it easy to interact with Bluetooth LE aka Bluetooth 4.0 using Go.

It is written using the [TinyGo Bluetooth](tinygo.org/x/bluetooth) package.

Learn more about Bluetooth LE at <http://en.wikipedia.org/wiki/Bluetooth_low_energy>

Drivers for several BLE Services can be found in the according [driver folder](https://github.com/hybridgroup/gobot/tree/release/drivers/ble).

## How to Install

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

### macOS

You need to have XCode installed to be able to compile code that uses the Gobot BLE adaptor on macOS. This is because the
`bluetooth` package uses a CGo based implementation.

### Ubuntu

Everything should already just compile on most Linux systems.

### Windows

You will need to have a GCC compiler such as [mingw-w64](https://github.com/mingw-w64/mingw-w64) installed in order to use
BLE on Windows.

## How To Connect

When using BLE a "peripheral" aka "server" is something you connect to such a a pulse meter. A "central" aka "client" is
what does the connecting, such as your computer or mobile phone.

You need to know the BLE ID of the peripheral you want to connect to. The Gobot BLE client adaptor also lets you connect
to a peripheral by friendly name.

### Connect on Ubuntu

On Linux the BLE code will need to run as a root user account. The easiest way to accomplish this is probably to use
`go build` to build your program, and then to run the requesting executable using `sudo`.

For example:

    go build examples/minidrone.go
    sudo ./minidrone AA:BB:CC:DD:EE

### Connect on Windows

Hopefully coming soon...

## How to Use

Here is an example that uses the BLE "Battery" service to retrieve the current change level of the peripheral device:

```go
package main

import (
  "fmt"
  "os"
  "time"

  "gobot.io/x/gobot/v2"
  "gobot.io/x/gobot/v2/drivers/ble"
  "gobot.io/x/gobot/v2/platforms/bleclient"
)

func main() {
  bleAdaptor := bleclient.NewAdaptor(os.Args[1])
  battery := ble.NewBatteryDriver(bleAdaptor)

  work := func() {
    gobot.Every(5*time.Second, func() {
      level, err := battery.GetBatteryLevel()
      if err != nil {
        fmt.Println(err)
      }
      fmt.Println("Battery level:", level)
    })
  }

  robot := gobot.NewRobot("bleBot",
    []gobot.Connection{bleAdaptor},
    []gobot.Device{battery},
    work,
  )

  if err := robot.Start(); err != nil {
		panic(err)
	}
}
```
