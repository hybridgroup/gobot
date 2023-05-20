# Jetson Nano

The Jetson Nano is ARM based single board computer with digital & PWM GPIO, and i2c interfaces built in.

The Gobot adaptor for the Jetson Nano should support Jetno Nano.

For more info about the Jetson Nano platform, click [here](https://developer.nvidia.com/embedded/jetson-nano/).

## How to Install

We recommend updating to the latest jetson-nano OS when using the Jetson Nano, however Gobot should also support older versions of the OS, should your application require this.

You would normally install Go and Gobot on your workstation. Once installed, cross compile your program on your workstation, transfer the final executable to your Jetson Nano, and run the program on the Jetson Nano as documented here.

```
go get -d -u gobot.io/x/gobot/v2/...
```

## How to Use

The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.

```go
package main

import (
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/jetson"
)

func main() {
	r := jetson.NewAdaptor()
	led := gpio.NewLedDriver(r, "40")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{r},
		[]gobot.Device{led},
		work,
	)

	robot.Start()

}
```

## How to Connect

### Compiling

Once you have compiled your code, you can upload your program and execute it on the Jetson Nano from your workstation using the `scp` and `ssh` commands like this:

```bash
$ scp jetson-nano_blink jn@192.168.1.xxx:/home/jn/
$ ssh -t jn@192.168.1.xxx "./jetson-nano_blink"
```
