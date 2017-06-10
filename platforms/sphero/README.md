# Sphero

Sphero is a sophisticated and programmable robot housed in a polycarbonate sphere shell.

The Gobot Sphero Adaptor & Driver makes it easy to interact with Sphero using Go. Once you have your Sphero setup and connected to your computer you can start writing code to make Sphero move, change direction, speed and colors, or detect Sphero events and execute some code when they occur.

Learn more about the Sphero robot go here: http://www.gosphero.com/

## How to Install
```
go get -d -u gobot.io/x/gobot/...
```

## How To Connect

### OSX

In order to allow Gobot running on your Mac to access the Sphero, go to "Bluetooth > Open Bluetooth Preferences > Sharing Setup" and make sure that "Bluetooth Sharing" is checked.

Now you must pair with the Sphero. Open System Preferences > Bluetooth. Now with the Bluetooth devices windows open,  smack the Sphero until it starts flashing three colors. You should see "Sphero-XXX" pop up as available devices where "XXX" is the first letter of the three colors the sphero is flashing. Pair with that device. Once paired your Sphero will be accessable through the serial device similarly named as `/dev/tty.Sphero-XXX-RN-SPP`

### Ubuntu

Connecting to the Sphero from Ubuntu or any other Linux-based OS can be done entirely from the command line using [Gort](https://gobot.io/x/gort) CLI commands. Here are the steps.

Find the address of the Sphero, by using:
```
gort scan bluetooth
```

Pair to Sphero using this command (substituting the actual address of your Sphero):
```
gort bluetooth pair <address>
```

Connect to the Sphero using this command (substituting the actual address of your Sphero):
```
gort bluetooth connect <address>
```

### Windows

You should be able to pair your Sphero using your normal system tray applet for Bluetooth, and then connect to the COM port that is bound to the device, such as `COM3`.

## How to Use

Example of a simple program that makes the Sphero roll.

```go
package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/sphero"
)

func main() {
	adaptor := sphero.NewAdaptor("/dev/rfcomm0")
	driver := sphero.NewSpheroDriver(adaptor)

	work := func() {
		gobot.Every(3*time.Second, func() {
			driver.Roll(30, uint16(gobot.Rand(360)))
		})
	}

	robot := gobot.NewRobot("sphero",
		[]gobot.Connection{adaptor},
		[]gobot.Device{driver},
		work,
	)

	robot.Start()
}
```
