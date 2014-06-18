# Sphero

This package provides the Gobot adaptor and driver for the [Sphero](http://www.gosphero.com/) robot from Orbotix .

## Installing
```
go get github.com/hybridgroup/gobot && go install github.com/hybridgroup/platforms/sphero
```

## How To Connect

### OSX

In order to allow Gobot running on your Mac to access the Sphero, go to "Bluetooth > Open Bluetooth Preferences > Sharing Setup" and make sure that "Bluetooth Sharing" is checked.

Now you must pair with the Sphero. Open System Preferences > Bluetooth. Now with the Bluetooth devices windows open,  smack the Sphero until it starts flashing three colors. You should see "Sphero-XXX" pop up as available devices where "XXX" is the first letter of the three colors the sphero is flashing. Pair with that device. Once paired your Sphero will be accessable through the serial device similarly named as `/dev/tty.Sphero-XXX-RN-SPP`

### Ubuntu

Connecting to the Sphero from Ubuntu or any other Linux-based OS can be done entirely from the command line using [Gort](https://github.com/hybridgroup/gort) CLI commands. Here are the steps.

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

## Example

```go
package main

import (
  "github.com/hybridgroup/gobot"
  "github.com/hybridgroup/gobot/platforms/sphero"
  "time"
)

func main() {
  gbot := gobot.NewGobot()

  adaptor := sphero.NewSpheroAdaptor("Sphero", "/dev/rfcomm0")
  ball := sphero.NewSpheroDriver(adaptor, "sphero")

  work := func() {
    gobot.Every(3*time.Second, func() {
      ball.Roll(30, uint16(gobot.Rand(360)))
    })
  }

  gbot.Robots = append(gbot.Robots,
    gobot.NewRobot("sphero", []gobot.Connection{adaptor}, []gobot.Device{ball}, work))

  gbot.Start()
}
```