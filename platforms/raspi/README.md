# Raspi

This package contains the Gobot adaptor for the [Raspberry Pi](http://www.raspberrypi.org/).

## Getting Started

First you must install the appropriate Go packages

```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/raspi
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
``` bash
$ scp raspi_blink pi@192.168.1.xxx:/home/pi/
```

and execute it on your Raspberry Pi with
```bash
$ ./raspi_blink
```

## Example

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
