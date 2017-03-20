// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/mavlink"
	common "gobot.io/x/gobot/platforms/mavlink/common"
)

func main() {
	adaptor := mavlink.NewAdaptor("/dev/ttyACM0")
	iris := mavlink.NewDriver(adaptor)

	work := func() {
		iris.Once(mavlink.PacketEvent, func(data interface{}) {
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

		iris.On(mavlink.MessageEvent, func(data interface{}) {
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

	robot.Start()
}
