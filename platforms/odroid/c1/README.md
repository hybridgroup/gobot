# ODroid C1

The ODroid-C1 is an inexpensive and popular ARM based single board computer with digital & PWM GPIO, and i2c interfaces built in.

The ODroid-C1 is a credit-card-sized single-board computer developed in South Korea.

For more info about the ODroid-C1 device, click [here](http://www.hardkernel.com/main/products/prdt_info.php?g_code=G141578608433/).

## How to Install

First you must install the appropriate Go packages

```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/odroid
```

#### Cross compiling for the ODroid C1
You must first configure your Go environment for linux cross compiling

```bash
$ cd $GOROOT/src
$ GOOS=linux GOARCH=arm ./make.bash --no-clean

```

Then compile your Gobot program with

```bash
$ GOARM=7 GOARCH=arm GOOS=linux go build examples/odroidc1_blink.go
```

Then you can simply upload your program over the network from your host computer to the ODroid

```bash
$ scp odroidc1_blink pi@192.168.1.xxx:/home/odroid/
```

and execute it on your ODroidC1 with

```bash
$ ./odroidc1_blink
```

## How to Use

```go
package main

import (
        "time"

        "github.com/hybridgroup/gobot"
        "github.com/hybridgroup/gobot/platforms/gpio"
        "github.com/hybridgroup/gobot/platforms/odroid/c1"
)

func main() {
        gbot := gobot.NewGobot()

        r := c1.NewODroidC1Adaptor("c1")
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
