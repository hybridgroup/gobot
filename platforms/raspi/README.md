# Raspberry Pi

The Raspberry Pi is an inexpensive and popular ARM based single board computer with digital & PWM GPIO, and i2c interfaces built in.

The Raspberry Pi is a credit-card-sized single-board computer developed in the UK by the Raspberry Pi Foundation with the intention of promoting the teaching of basic computer science in schools

For more info about the Raspberry Pi platform, click [here](http://www.raspberrypi.org/).

## How to Install

```
go get -d -u gobot.io/x/gobot/... && go install gobot.io/x/gobot/platforms/raspi
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

Simply compile your Gobot program like this:

```bash
$ GOARM=6 GOARCH=arm GOOS=linux examples/raspi_blink.go
```

Then you can simply upload your program over the network from your host computer to the Raspi

```bash
$ scp raspi_blink pi@192.168.1.xxx:/home/pi/
```

and execute it on your Raspberry Pi with

```bash
$ ./raspi_blink
```

### Enabling PWM output on GPIO pins.

For extended PWM support on the Raspberry Pi, you will need to use a program called pi-blaster. You can follow the instructions for pi-blaster install in the pi-blaster repo here:

[https://github.com/sarfata/pi-blaster](https://github.com/sarfata/pi-blaster)

### Special note for Raspian Wheezy users

The Golang version installed from the default package repositories is very old and will not compile Gobot. You can install go 1.4 as follows:

```bash
$ wget -O - http://dave.cheney.net/paste/go1.4.linux-arm~multiarch-armv6-1.tar.gz|sudo tar -xzC /usr/local -f -

$ echo '# Setup for golang' |sudo tee /etc/profile.d/golang.sh
$ echo 'PATH=$PATH:/usr/local/go/bin'|sudo tee -a /etc/profile.d/golang.sh

$ source /etc/profile.d/golang.sh
```
