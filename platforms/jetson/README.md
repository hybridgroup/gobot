# Jetson Nano

The Jetson Nano is ARM based single board computer with digital & PWM GPIO, and i2c interfaces built in.

The Gobot adaptor for the Jetson Nano should support Jetno Nano.

For more info about the Jetson Nano platform, click [here](https://developer.nvidia.com/embedded/jetson-nano/).

## How to Install

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

We recommend updating to the latest jetson-nano OS when using the Jetson Nano, however Gobot should also support older
versions of the OS, should your application require this.

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
      if err := led.Toggle(); err != nil {
				fmt.Println(err)
			}
    })
  }

  robot := gobot.NewRobot("blinkBot",
    []gobot.Connection{r},
    []gobot.Device{led},
    work,
  )

  if err := robot.Start(); err != nil {
		panic(err)
	}

}
```

## How to Connect

### Compiling

Once you have compiled your code, you can upload your program and execute it on the Jetson Nano from your workstation using
the `scp` and `ssh` commands like this:

```sh
scp jetson-nano_blink jn@192.168.1.xxx:/home/jn/
ssh -t jn@192.168.1.xxx "./jetson-nano_blink"
```
