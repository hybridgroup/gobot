# Joystick

You can use Gobot with many USB joysticks and game controllers.

Current configurations included:

- Dualshock3 game controller
- Dualshock4 game controller
- Dualsense game controller
- Thrustmaster T-Flight Hotas X Joystick
- XBox360 game controller
- XBox360 "Rock Band" drum controller
- Nintendo Switch Joy-Con controller pair

## How to Install

Any platform specific info here...

### macOS


### Linux (Ubuntu and Raspbian)


### Windows


## How to Use

Controller configurations are stored in Gobot, but you can also use external file in JSON format. Take a look at the `configs` directory for examples.

## How to Connect

Plug your USB joystick or game controller into your USB port. If your device is supported by your operating system, it might prompt you to install some system drivers.

For the Dualshock4, you must pair the device with your computers Bluetooth interface first, before running your Gobot program.

## Examples

This small program receives joystick and button press events from an PlayStation 3 game controller.

```go
package main

import (
  "fmt"

  "gobot.io/x/gobot/v2"
  "gobot.io/x/gobot/v2/platforms/joystick"
)

func main() {
  joystickAdaptor := joystick.NewAdaptor("0")
  stick := joystick.NewDriver(joystickAdaptor, "dualshock3",
  )

  work := func() {
    // buttons
    _ = stick.On(joystick.SquarePress, func(data interface{}) {
      fmt.Println("square_press")
    })
    _ = stick.On(joystick.SquareRelease, func(data interface{}) {
      fmt.Println("square_release")
    })
    _ = stick.On(joystick.TrianglePress, func(data interface{}) {
      fmt.Println("triangle_press")
    })
    _ = stick.On(joystick.TriangleRelease, func(data interface{}) {
      fmt.Println("triangle_release")
    })
    _ = stick.On(joystick.CirclePress, func(data interface{}) {
      fmt.Println("circle_press")
    })
    _ = stick.On(joystick.CircleRelease, func(data interface{}) {
      fmt.Println("circle_release")
    })
    _ = stick.On(joystick.XPress, func(data interface{}) {
      fmt.Println("x_press")
    })
    _ = stick.On(joystick.XRelease, func(data interface{}) {
      fmt.Println("x_release")
    })
    _ = stick.On(joystick.StartPress, func(data interface{}) {
      fmt.Println("start_press")
    })
    _ = stick.On(joystick.StartRelease, func(data interface{}) {
      fmt.Println("start_release")
    })
    _ = stick.On(joystick.SelectPress, func(data interface{}) {
      fmt.Println("select_press")
    })
    _ = stick.On(joystick.SelectRelease, func(data interface{}) {
      fmt.Println("select_release")
    })

    // joysticks
    _ = stick.On(joystick.LeftX, func(data interface{}) {
      fmt.Println("left_x", data)
    })
    _ = stick.On(joystick.LeftY, func(data interface{}) {
      fmt.Println("left_y", data)
    })
    _ = stick.On(joystick.RightX, func(data interface{}) {
      fmt.Println("right_x", data)
    })
    _ = stick.On(joystick.RightY, func(data interface{}) {
      fmt.Println("right_y", data)
    })

    // triggers
    _ = stick.On(joystick.R1Press, func(data interface{}) {
      fmt.Println("R1Press", data)
    })
    _ = stick.On(joystick.R2Press, func(data interface{}) {
      fmt.Println("R2Press", data)
    })
    _ = stick.On(joystick.L1Press, func(data interface{}) {
      fmt.Println("L1Press", data)
    })
    _ = stick.On(joystick.L2Press, func(data interface{}) {
      fmt.Println("L2Press", data)
    })
  }

  robot := gobot.NewRobot("joystickBot",
    []gobot.Connection{joystickAdaptor},
    []gobot.Device{stick},
    work,
  )

  if err := robot.Start(); err != nil {
		panic(err)
	}
}
```

## How to Add A New Joystick

You can create a file similar to `joystick_dualshock3.go` and submit a pull request with the new configuration so others can use it as well.
