# Raspi

The Raspberry Pi is an inexpensive and popular ARM based single board computer with digital & PWM GPIO, and i2c interfaces built in.

The Raspberry Pi is a credit-card-sized single-board computer developed in the UK by the Raspberry Pi Foundation with the intention of promoting the teaching of basic computer science in schools

For more info about the Raspberry Pi platform, click [here](http://www.raspberrypi.org/).

## How to Install

First you must install the appropriate Go packages

```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/raspi
```

### Enabling PWM output on GPIO pins.

You need to install and have pi-blaster running in the raspberry-pi, you can follow the instructions for pi-blaster install in the pi-blaster repo here:

[https://github.com/sarfata/pi-blaster](https://github.com/sarfata/pi-blaster)

### Special note for Raspian Wheezy users

The go vesion installed from the default package repositories is very old and will not compile gobot. You can install go 1.4 as follows:

```bash
$ wget -O - http://dave.cheney.net/paste/go1.4.linux-arm~multiarch-armv6-1.tar.gz|sudo tar -xzC /usr/local -f -

$ echo '# Setup for golang' |sudo tee /etc/profile.d/golang.sh
$ echo 'PATH=$PATH:/usr/local/go/bin'|sudo tee -a /etc/profile.d/golang.sh

$ source /etc/profile.d/golang.sh
```

#### Cross compiling for the Raspberry Pi
You must first configure your Go environment for linux cross compiling

```bash
$ cd $GOROOT/src
$ GOOS=linux GOARCH=arm ./make.bash --no-clean

```

Then compile your Gobot program with

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

## How to Use

```go
package main

import (
        "time"

        "github.com/hybridgroup/gobot"
        "github.com/hybridgroup/gobot/platforms/gpio"
        "github.com/hybridgroup/gobot/platforms/raspi"
)

func main() {
        gbot := gobot.NewGobot()

        r := raspi.NewRaspiAdaptor("raspi")
        led := gpio.NewLedDriver(r, "led", "7")

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

        gbot.AddRobot(robot)

        gbot.Start()
}
```
