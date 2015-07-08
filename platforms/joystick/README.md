# Joystick

You can use Gobot with a Dualshock3 game controller, an XBox360 game controller, or any other USB joystick or game controller that is compatible with [Simple DirectMedia Layer](http://www.libsdl.org/).

## How to Install

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
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/joystick
```

## How to Use

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

## How to Connect

Plug your USB joystick or game controller into your USB port. If your device is supported by SDL, you are now ready.

## Examples

This small program receives joystick and button press events from an PlayStation 3 game controller.

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
	joystick := joystick.NewJoystickDriver(joystickAdaptor,
		"ps3",
		"./platforms/joystick/configs/dualshock3.json",
	)

	work := func() {
		gobot.On(joystick.Event("square_press"), func(data interface{}) {
			fmt.Println("square_press")
		})
		gobot.On(joystick.Event("square_release"), func(data interface{}) {
			fmt.Println("square_release")
		})
		gobot.On(joystick.Event("triangle_press"), func(data interface{}) {
			fmt.Println("triangle_press")
		})
		gobot.On(joystick.Event("triangle_release"), func(data interface{}) {
			fmt.Println("triangle_release")
		})
		gobot.On(joystick.Event("left_x"), func(data interface{}) {
			fmt.Println("left_x", data)
		})
		gobot.On(joystick.Event("left_y"), func(data interface{}) {
			fmt.Println("left_y", data)
		})
		gobot.On(joystick.Event("right_x"), func(data interface{}) {
			fmt.Println("right_x", data)
		})
		gobot.On(joystick.Event("right_y"), func(data interface{}) {
			fmt.Println("right_y", data)
		})
	}

	robot := gobot.NewRobot("joystickBot",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{joystick},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
```