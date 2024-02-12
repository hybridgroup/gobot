# MegaPi

The [MegaPi](http://learn.makeblock.com/en/megapi/) is a motor controller by MakeBlock that is compatible with the
Raspberry Pi.

The code is based on a python implementation that can be found [here](https://github.com/Makeblock-official/PythonForMegaPi).

## How to Install

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

## How to Use

```go
package main

import (
  "fmt"
  "time"

  "gobot.io/x/gobot/v2"
  "gobot.io/x/gobot/v2/drivers/serial/megapi"
  "gobot.io/x/gobot/v2/platforms/serialport"
)

func main() {
  // use "/dev/ttyUSB0" if connecting with USB cable
  // use "/dev/ttyAMA0" on devices older than Raspberry Pi 3 Model B
  adaptor := serialport.NewAdaptor("/dev/ttyS0", serialport.WithName("MegaPi"))
  motor := megapi.NewMotorDriver(adaptor, 1)

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
    []gobot.Connection{adaptor},
    []gobot.Device{motor},
    work,
  )

  if err := robot.Start(); err != nil {
    panic(err)
  }
}
```
