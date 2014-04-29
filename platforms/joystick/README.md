# Gobot Adaptor for Joysticks and Game Controllers

Gobot (http://gobot.io/) is a library for robotics and physical computing using Go

This repository contains the Gobot adaptor for the PS3 controller, Xbox 360 controller, or any other joysticks and game controllers that are compatible with Simple DirectMedia Layer (SDL) (http://www.libsdl.org/).

For more information about Gobot, check out the github repo at
https://github.com/hybridgroup/gobot

## Getting Started

Installing gobot-joystick requires `sdl2` to be installed on your system

### Ubuntu

```
$ sudo apt-get install libsdl2-2.0-0
```

Now you can install `gobot-joystick` with
```
$ go get github.com/hybridgroup/gobot-joystick
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
	"github.com/hybridgroup/gobot-joystick"
)

func main() {
	joystickAdaptor := new(gobotJoystick.JoystickAdaptor)
	joystickAdaptor.Name = "ps3"
	joystickAdaptor.Params = map[string]interface{}{
		"config": "./configs/dualshock3.json",
	}

	joystick := gobotJoystick.NewJoystick(joystickAdaptor)
	joystick.Name = "ps3"

	work := func() {
		gobot.On(joystick.Events["square_press"], func(data interface{}) {
			fmt.Println("square_press")
		})
		gobot.On(joystick.Events["square_release"], func(data interface{}) {
			fmt.Println("square_release")
		})
		gobot.On(joystick.Events["triangle_press"], func(data interface{}) {
			fmt.Println("triangle_press")
		})
		gobot.On(joystick.Events["triangle_release"], func(data interface{}) {
			fmt.Println("triangle_release")
		})
		gobot.On(joystick.Events["left_x"], func(data interface{}) {
			fmt.Println("left_x", data)
		})
		gobot.On(joystick.Events["left_y"], func(data interface{}) {
			fmt.Println("left_y", data)
		})
		gobot.On(joystick.Events["right_x"], func(data interface{}) {
			fmt.Println("right_x", data)
		})
		gobot.On(joystick.Events["right_y"], func(data interface{}) {
			fmt.Println("right_y", data)
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{joystickAdaptor},
		Devices:     []gobot.Device{joystick},
		Work:        work,
	}

	robot.Start()
}
```

## Documentation
We're busy adding documentation to our web site at http://gobot.io/ please check there as we continue to work on Gobot

Thank you!

## Contributing

* All patches must be provided under the Apache 2.0 License
* Please use the -s option in git to "sign off" that the commit is your work and you are providing it under the Apache 2.0 License
* Submit a Github Pull Request to the appropriate branch and ideally discuss the changes with us in IRC.
* We will look at the patch, test it out, and give you feedback.
* Avoid doing minor whitespace changes, renamings, etc. along with merged content. These will be done by the maintainers from time to time but they can complicate merges and should be done seperately.
* Take care to maintain the existing coding style.
* Add unit tests for any new or changed functionality.
* All pull requests should be "fast forward"
  * If there are commits after yours use “git rebase -i <new_head_branch>”
  * If you have local changes you may need to use “git stash”
  * For git help see [progit](http://git-scm.com/book) which is an awesome (and free) book on git


## License
Copyright (c) 2013-2014 The Hybrid Group. Licensed under the Apache 2.0 license.