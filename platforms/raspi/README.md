# Raspberry Pi

The Raspberry Pi is an inexpensive and popular ARM based single board computer with digital & PWM GPIO, and i2c interfaces built in.

The Gobot adaptor for the Raspberry Pi should support all of the various Raspberry Pi boards such as the Raspberry Pi 3 Model B, Raspberry Pi 2 Model B, Raspberry Pi 1 Model A+, Raspberry Pi Zero, and Raspberry Pi Zero W.

For more info about the Raspberry Pi platform, click [here](http://www.raspberrypi.org/).

## How to Install

We recommend updating to the latest Raspian Jessie OS when using the Raspberry Pi, however Gobot should also support older versions of the OS, should your application require this.

You would normally install Go and Gobot on your workstation. Once installed, cross compile your program on your workstation, transfer the final executable to your Raspberry Pi, and run the program on the Raspberry Pi as documented here.

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use

The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.

```go
package main

import (
        "time"

        "gobot.io/x/gobot"
        "gobot.io/x/gobot/drivers/gpio"
        "gobot.io/x/gobot/platforms/raspi"
)

func main() {
        r := raspi.NewAdaptor()
        led := gpio.NewLedDriver(r, "7")

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

Compile your Gobot program on your workstation like this:

```bash
$ GOARM=6 GOARCH=arm GOOS=linux go build examples/raspi_blink.go
```

Use the following `GOARM` values to compile depending on which model Raspberry Pi you are using:

`GOARM=6` (Raspberry Pi A, A+, B, B+, Zero)
`GOARM=7` (Raspberry Pi 2, 3)

Once you have compiled your code, you can upload your program and execute it on the Raspberry Pi from your workstation using the `scp` and `ssh` commands like this:

```bash
$ scp raspi_blink pi@192.168.1.xxx:/home/pi/
$ ssh -t pi@192.168.1.xxx "./raspi_blink"
```

### Enabling PWM output on GPIO pins.

For extended PWM support on the Raspberry Pi, you will need to use a program called pi-blaster. You can follow the instructions for pi-blaster install in the pi-blaster repo here:

[https://github.com/sarfata/pi-blaster](https://github.com/sarfata/pi-blaster)
