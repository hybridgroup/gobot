# C.H.I.P.

The [C.H.I.P.](http://www.getchip.com/) is a small, inexpensive ARM based single board computer, with many different IO interfaces available on the [pin headers](http://docs.getchip.com/#pin-headers).

We recommend updating to the latest Debian OS when using the C.H.I.P., however Gobot should also support older versions of the OS, should your application require this.

For documentation about the C.H.I.P. platform click [here](http://docs.getchip.com/).

## How to Install
```
go get -d -u gobot.io/x/gobot/... && go install gobot.io/x/gobot/platforms/chip
```

## How to Use

The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.

```go
package main

import (
    "fmt"

    "gobot.io/x/gobot"
    "gobot.io/x/gobot/drivers/gpio"
    "gobot.io/x/gobot/platforms/chip"
)

func main() {
    chipAdaptor := chip.NewAdaptor()
    button := gpio.NewButtonDriver(chipAdaptor, "XIO-P0")

    work := func() {
        gobot.On(button.Event("push"), func(data interface{}) {
            fmt.Println("button pressed")
        })

        gobot.On(button.Event("release"), func(data interface{}) {
            fmt.Println("button released")
        })
    }

    robot := gobot.NewRobot("buttonBot",
        []gobot.Connection{chipAdaptor},
        []gobot.Device{button},
        work,
    )

    robot.Start()
}
```

## How to Connect

### Compiling

Compile your Gobot program like this:

```bash
$ GOARM=7 GOARCH=arm GOOS=linux go build examples/chip_button.go
```

Then you can simply upload your program to the CHIP and execute it with

```bash
$ scp chip_button root@192.168.1.xx:
$ ssh -t root@192.168.1.xx "./chip_button"
```
