# Joystick

You can use Gobot with any USB joystick or game controller that is compatible with [Simple DirectMedia Layer](http://www.libsdl.org/).

Current configurations included:
- Dualshock3 game controller
- Dualshock4 game controller
- Thrustmaster T-Flight Hotas X Joystick
- XBox360 game controller
- XBox360 "Rock Band" drum controller

## How to Install

This package requires `sdl2` to be installed on your system

### OSX

To install `sdl2` on OSX using Homebrew:

```
$ brew install sdl2
```

To use an XBox360 controller on OS X, you will most likely need to install additional software such as [https://github.com/360Controller/360Controller](https://github.com/360Controller/360Controller).

### Ubuntu

```
wget https://www.libsdl.org/release/SDL2-2.0.8.tar.gz
tar -zxvf SDL2-2.0.8.tar.gz
cd SDL2-2.0.8/
./configure && make && sudo make install
```

Now you can install the package with

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use

Controller configurations are stored in Gobot it, but you can also use external file in JSON format. Take a look at the `configs` directory for examples.


## How to Connect

Plug your USB joystick or game controller into your USB port. If your device is supported by SDL, you are now ready.

For the Dualshock4, you must pair the device with your computers Bluetooth interface first, before running your Gobot program.

## Examples

This small program receives joystick and button press events from an PlayStation 3 game controller.

```go
package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/joystick"
)

func main() {
	joystickAdaptor := joystick.NewAdaptor()
	stick := joystick.NewDriver(joystickAdaptor, "dualshock3",
	)

	work := func() {
		// buttons
		stick.On(joystick.SquarePress, func(data interface{}) {
			fmt.Println("square_press")
		})
		stick.On(joystick.SquareRelease, func(data interface{}) {
			fmt.Println("square_release")
		})
		stick.On(joystick.TrianglePress, func(data interface{}) {
			fmt.Println("triangle_press")
		})
		stick.On(joystick.TriangleRelease, func(data interface{}) {
			fmt.Println("triangle_release")
		})
		stick.On(joystick.CirclePress, func(data interface{}) {
			fmt.Println("circle_press")
		})
		stick.On(joystick.CircleRelease, func(data interface{}) {
			fmt.Println("circle_release")
		})
		stick.On(joystick.XPress, func(data interface{}) {
			fmt.Println("x_press")
		})
		stick.On(joystick.XRelease, func(data interface{}) {
			fmt.Println("x_release")
		})
		stick.On(joystick.StartPress, func(data interface{}) {
			fmt.Println("start_press")
		})
		stick.On(joystick.StartRelease, func(data interface{}) {
			fmt.Println("start_release")
		})
		stick.On(joystick.SelectPress, func(data interface{}) {
			fmt.Println("select_press")
		})
		stick.On(joystick.SelectRelease, func(data interface{}) {
			fmt.Println("select_release")
		})

		// joysticks
		stick.On(joystick.LeftX, func(data interface{}) {
			fmt.Println("left_x", data)
		})
		stick.On(joystick.LeftY, func(data interface{}) {
			fmt.Println("left_y", data)
		})
		stick.On(joystick.RightX, func(data interface{}) {
			fmt.Println("right_x", data)
		})
		stick.On(joystick.RightY, func(data interface{}) {
			fmt.Println("right_y", data)
		})

		// triggers
		stick.On(joystick.R1Press, func(data interface{}) {
			fmt.Println("R1Press", data)
		})
		stick.On(joystick.R2Press, func(data interface{}) {
			fmt.Println("R2Press", data)
		})
		stick.On(joystick.L1Press, func(data interface{}) {
			fmt.Println("L1Press", data)
		})
		stick.On(joystick.L2Press, func(data interface{}) {
			fmt.Println("L2Press", data)
		})
	}

	robot := gobot.NewRobot("joystickBot",
		[]gobot.Connection{joystickAdaptor},
		[]gobot.Device{stick},
		work,
	)

	robot.Start()
}
```

## How to Add A New Joystick

In the `bin` directory for this package is a CLI utility program that scans for SDL joystick events, and displays the ID and value:

```
$ go run ./platforms/joystick/bin/scanner.go 
Joystick 0 connected
[6625 ms] Axis: 1       value:-22686
[6641 ms] Axis: 1       value:-32768
[6836 ms] Axis: 1       value:-18317
[6852 ms] Axis: 1       value:0
[8663 ms] Axis: 3       value:-32768
[8873 ms] Axis: 3       value:0
[10183 ms] Axis: 0      value:-24703
[10183 ms] Axis: 0      value:-32768
[10313 ms] Axis: 1      value:-3193
[10329 ms] Axis: 1      value:0
[10345 ms] Axis: 0      value:0
```

You can use the output from this program to create a JSON file for the various buttons and axes on your joystick/gamepad. You could also create a file similar to `joystick_dualshock3.go` and submit a pull request with the new configuration so others can use it as well.
