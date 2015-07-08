# Mavlink

For information on the MAVlink communication protocol click [here](http://qgroundcontrol.org/mavlink/start).

## How to Install

```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/mavlink

```

## How to Use

```go
package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/mavlink"
	common "github.com/hybridgroup/gobot/platforms/mavlink/common"
)

func main() {
	gbot := gobot.NewGobot()

	adaptor := mavlink.NewMavlinkAdaptor("iris", "/dev/ttyACM0")
	iris := mavlink.NewMavlinkDriver(adaptor, "iris")

	work := func() {
		gobot.Once(iris.Event("packet"), func(data interface{}) {
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

		gobot.On(iris.Event("message"), func(data interface{}) {
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

	gbot.AddRobot(robot)

	gbot.Start()
}
```
