/*
Package mavlink contains the Gobot adaptor and driver for the MAVlink Communication Protocol.

Installing:

	Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

Example:

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

For further information refer to mavlink README:
https://github.com/hybridgroup/gobot/blob/release/platforms/mavlink/README.md
*/
package mavlink // import "gobot.io/x/gobot/v2/platforms/mavlink"
