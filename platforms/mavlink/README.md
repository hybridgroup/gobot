# Mavlink

For information on the MAVlink communication protocol click [here](https://mavlink.io/).

This package supports Mavlink over serial (such as a
[SiK modem](http://ardupilot.org/copter/docs/common-sik-telemetry-radio.html))
and Mavlink over UDP (such as via
[mavproxy](https://github.com/ArduPilot/MAVProxy)).  Serial is useful
for point to point links, and UDP is useful for where you have
multiple simultaneous clients such as the robot and
[QGroundControl](http://qgroundcontrol.com/).

As at 2018-04, this package supports Mavlink 1.0 only.  If the robot
doesn't receiving data then check that the other devices are
configured to send version 1.0 frames.

For Mavlink 2.0 please refer to [gomavlib](https://github.com/bluenviron/gomavlib).

## How to Install

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

## How to Use

```go
package main

import (
  "fmt"

  "gobot.io/x/gobot/v2"
  "gobot.io/x/gobot/v2/platforms/mavlink"
  common "gobot.io/x/gobot/v2/platforms/mavlink/common"
)

func main() {
  adaptor := mavlink.NewAdaptor("/dev/ttyACM0")
  iris := mavlink.NewDriver(adaptor)

  work := func() {
    iris.Once(iris.Event("packet"), func(data interface{}) {
      packet := data.(*common.MAVLinkPacket)

      dataStream := common.NewRequestDataStream(100,
        packet.SystemID,
        packet.ComponentID,
        4,
        1,
      )
      iris.SendPacket(common.CraftMAVLinkPacket(packet.SystemID,
        packet.ComponentID,
        dataStream,
      ))
    })

    iris.On(iris.Event("message"), func(data interface{}) {
      if data.(common.MAVLinkMessage).Id() == 30 {
        message := data.(*common.Attitude)
        fmt.Println("Attitude")
        fmt.Println("TIME_BOOT_MS", message.TIME_BOOT_MS)
        fmt.Println("ROLL", message.ROLL)
        fmt.Println("PITCH", message.PITCH)
        fmt.Println("YAW", message.YAW)
        fmt.Println("ROLLSPEED", message.ROLLSPEED)
        fmt.Println("PITCHSPEED", message.PITCHSPEED)
        fmt.Println("YAWSPEED", message.YAWSPEED)
        fmt.Println("")
      }
    })
  }

  robot := gobot.NewRobot("mavBot",
    []gobot.Connection{adaptor},
    []gobot.Device{iris},
    work,
  )

  if err := robot.Start(); err != nil {
		panic(err)
	}
}
```

## How to use: UDP

``` go
  adaptor := mavlink.NewUDPAdaptor(":14550")
```

To test, install Mavproxy and set it up to listen on serial and repeat
over UDP:

`$ mavproxy.py --out=udpbcast:192.168.0.255:14550`

Change the address to the broadcast address of your subnet.
