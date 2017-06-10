# MegaPi

The [MegaPi](http://learn.makeblock.com/en/megapi/) is a motor controller by MakeBlock that is compatible with the Raspberry Pi.

The code is based on a python implementation that can be found [here](https://github.com/Makeblock-official/PythonForMegaPi).

## How to Install

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use

```go
package main

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/megapi"
	"time"
)

func main() {
	// use "/dev/ttyUSB0" if connecting with USB cable
	// use "/dev/ttyAMA0" on devices older than Raspberry Pi 3 Model B
	megaPiAdaptor := megapi.NewAdaptor("/dev/ttyS0")
	motor := megapi.NewMotorDriver(megaPiAdaptor, 1)

	work := func() {
		speed := int16(0)
		fadeAmount := int16(30)

		gobot.Every(100*time.Millisecond, func() {
			motor.Speed(speed)
			speed = speed + fadeAmount
			if speed == 0 || speed == 300 {
				fadeAmount = -fadeAmount
			}
		})
	}

	robot := gobot.NewRobot("megaPiBot",
		[]gobot.Connection{megaPiAdaptor},
		[]gobot.Device{motor},
		work,
	)

	robot.Start()
}
```
