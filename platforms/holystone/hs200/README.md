# Holystone HS200

This package contains the Gobot driver for the Holystone HS200 drone.

For more information on this drone, go to:
http://www.holystone.com/product/Holy_Stone_HS200W_FPV_Drone_with_720P_HD_Live_Video_Wifi_Camera_2_4GHz_4CH_6_Axis_Gyro_RC_Quadcopter_with_Altitude_Hold,_Gravity_Sensor_and_Headless_Mode_Function_RTF,_Color_Red-39.html

## How to Install

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use
- Connect to the drone's Wi-Fi network and identify the drone/gateway IP address.
- Use that IP address when you create a new driver.
- Some drones appear to use a different TCP port (8080 vs. 8888?).  If the example doesn't work scan the drone for open ports or modify the driver not to use TCP.

Here is a sample of how you initialize and use the driver:

```go
package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/holystone/hs200"
)

func main() {
	drone := hs200.NewDriver("172.16.10.1:8888", "172.16.10.1:8080")

	work := func() {
		drone.TakeOff()

		gobot.After(5*time.Second, func() {
			drone.Land()
		})
	}

	robot := gobot.NewRobot("hs200",
		[]gobot.Connection{},
		[]gobot.Device{drone},
		work,
	)

	robot.Start()
}
```

## References
https://hackaday.io/project/19356/logs

https://github.com/lancecaraccioli/holystone-hs110w

## Random notes
- The hs200 sends out an RTSP video feed from its own board camera.  Not clear how this is turned on.  The data is apparently streamed over UDP. (Reference mentions rtsp://192.168.0.1/0 in VLC, I didn't try it!)
- The Android control app seems to be sending out the following TCP bytes for an unknown purpose:
`00 01 02 03 04 05 06 07 08 09 25 25` but the drone flies without a TCP connection.
- The drone apparently always replies "noact\r\n" over TCP.
- The app occasionally sends out 29 bytes long UDP packets besides the 11 byte control packet for an unknown purpose:
`26 e1 07 00 00 07 00 00 00 10 00 00 00 00 00 00 00 14 00 00 00 0e 00 00 00 03 00 00 00`
- The doesn't seem to be any telemetry coming out of the drone besides the video feed.
- The drone can sometimes be a little flaky.  Ensure you've got a fully charged battery, minimal Wi-Fi interference, various connectors on the drone all well seated.
- It's not clear whether the drone's remote uses Wi-Fi or not, possibly Wi-Fi is only for the mobile app.
