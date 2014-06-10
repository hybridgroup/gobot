# Joystick

This package provides the Gobot adaptor and drivers for the PS3 controller, Xbox 360 controller, or any other joysticks and game controllers that are compatible with [Simple DirectMedia Layer](http://www.libsdl.org/).

## Getting Started

This package requires `sdl2` to be installed on your system

### OSX

To install `sdl2` on OSX using Homebrew:

```
$ brew install sdl2
```

### Ubuntu

```
$ sudo apt-get install libsdl2-2.0-0
```

Now you can install the package with 
```
go get github.com/hybridgroup/gobot && go install github.com/hybridgroup/platforms/joystick
```

## Usage

Controller configurations are stored in JSON format. Here's an example configuration file for the Dualshock 3 controller
```json
{
    "name": "Sony PLAYSTATION(R)3 Controller",
    "guid": "030000004c0500006802000011010000",
    "axis": [
        {
            "name": "left_x",
            "id": 0
        },
        {
            "name": "left_y",
            "id": 1
        },
        {
            "name": "right_x",
            "id": 2
        },
        {
            "name": "right_y",
            "id": 3
        }
    ],
    "buttons": [
        {
            "name": "square",
            "id": 15
        },
        {
            "name": "triangle",
            "id": 12
        },
        {
            "name": "circle",
            "id": 13
        },
        {
            "name": "x",
            "id": 14
        },
        {
            "name": "up",
            "id": 4
        },
        {
            "name": "down",
            "id": 6
        },
        {
            "name": "left",
            "id": 7
        },
        {
            "name": "right",
            "id": 5
        },
        {
            "name": "left_stick",
            "id": 1
        },
        {
            "name": "right_stick",
            "id": 2
        },
        {
            "name": "l1",
            "id": 10
        },
        {
            "name": "l2",
            "id": 8
        },
        {
            "name": "r1",
            "id": 11
        },
        {
            "name": "r2",
            "id": 9
        },
        {
            "name": "start",
            "id": 3
        },
        {
            "name": "select",
            "id": 0
        },
        {
            "name": "home",
            "id": 16
        }
    ]
}
```

## Examples
```go
package main

import (
    "fmt"
    "github.com/hybridgroup/gobot"
    "github.com/hybridgroup/gobot/platforms/joystick"
)

func main() {
    gbot := gobot.NewGobot()
    joystickAdaptor := joystick.NewJoystickAdaptor("ps3")
    joystickDriver := joystick.NewJoystickDriver(joystickAdaptor, "ps3", "./platforms/joystick/configs/dualshock3.json")

    work := func() {
        gobot.On(joystickDriver.Events["square_press"], func(data interface{}) {
            fmt.Println("square_press")
        })
        gobot.On(joystickDriver.Events["square_release"], func(data interface{}) {
            fmt.Println("square_release")
        })
        gobot.On(joystickDriver.Events["triangle_press"], func(data interface{}) {
            fmt.Println("triangle_press")
        })
        gobot.On(joystickDriver.Events["triangle_release"], func(data interface{}) {
            fmt.Println("triangle_release")
        })
        gobot.On(joystickDriver.Events["left_x"], func(data interface{}) {
            fmt.Println("left_x", data)
        })
        gobot.On(joystickDriver.Events["left_y"], func(data interface{}) {
            fmt.Println("left_y", data)
        })
        gobot.On(joystickDriver.Events["right_x"], func(data interface{}) {
            fmt.Println("right_x", data)
        })
        gobot.On(joystickDriver.Events["right_y"], func(data interface{}) {
            fmt.Println("right_y", data)
        })
    }

    gbot.Robots = append(gbot.Robots,
        gobot.NewRobot("joystickBot", []gobot.Connection{joystickAdaptor}, []gobot.Device{joystickDriver}, work))

    gbot.Start()
}
```