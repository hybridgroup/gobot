# CHIP

The [CHIP](http://www.getchip.com/) is a small, inexpensive ARM based single board computer, with many different IO interfaces available on the [pin headers](http://docs.getchip.com/#pin-headers).

For documentation about the CHIP platform click [here](http://docs.getchip.com/).

## How to Install
```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/chip
```

## Cross compiling for the CHIP
If you're using Go version earlier than 1.5, you must first configure your Go environment for ARM linux cross compiling.

```bash
$ cd $GOROOT/src
$ GOOS=linux GOARCH=arm ./make.bash --no-clean
```

The above step is not required for Go >= 1.5

Then compile your Gobot program with

```bash
$ GOARM=7 GOARCH=arm GOOS=linux go build examples/chip_button.go
```

Then you can simply upload your program to the CHIP and execute it with

```bash
$ scp chip_button root@192.168.1.xx:
$ ssh -t root@192.168.1.xx "./chip_button"
```

## How to Use

```go
package main

import (
    "fmt"

    "github.com/hybridgroup/gobot"
    "github.com/hybridgroup/gobot/platforms/chip"
    "github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
    gbot := gobot.NewGobot()

    chipAdaptor := chip.NewChipAdaptor("chip")
    button := gpio.NewButtonDriver(chipAdaptor, "button", "XIO-P0")

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
    gbot.AddRobot(robot)
    gbot.Start()
}
```
