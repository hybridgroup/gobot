# DragonBoard™ 410c

The [DragonBoard 410c](http://www.96boards.org/product/dragonboard410c/), a product of Arrow Electronics, is the development board based on the mid-tier Qualcomm® Snapdragon™ 410E processor. It features advanced processing power, Wi-Fi, Bluetooth connectivity, and GPS, all packed into a board the size of a credit card.

Make sure you are using the latest linaro debian image.

## How to Install
```
go get -d -u gobot.io/x/gobot/... && go install gobot.io/x/gobot/platforms/dragonboard
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
$ GOARCH=arm64 GOOS=linux go build examples/chip_button.go
```

Then you can simply upload your program to the CHIP and execute it with

```bash
$ scp chip_button root@192.168.1.xx:
$ ssh -t root@192.168.1.xx "./chip_button"
```
