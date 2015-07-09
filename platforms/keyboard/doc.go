/*
Packge keyboard contains the Gobot drivers for keyboard support.

Installing:

Then you can install the package with:

	go get github.com/hybridgroup/gobot && go install github.com/hybridgroup/gobot/platforms/keyboard

Example:

	package main

	import (
		"fmt"

		"github.com/hybridgroup/gobot"
		"github.com/hybridgroup/gobot/platforms/keyboard"
	)

	func main() {
		gbot := gobot.NewGobot()

		keys := keyboard.NewKeyboardDriver("keyboard")

		work := func() {
			gobot.On(keys.Event("key"), func(data interface{}) {
				key := data.(keyboard.KeyEvent)

				if key.Key == keyboard.A {
					fmt.Println("A pressed!")
				} else {
					fmt.Println("keyboard event!", key, key.Char)
				}
			})
		}

		robot := gobot.NewRobot("keyboardbot",
			[]gobot.Connection{},
			[]gobot.Device{keys},
			work,
		)

		gbot.AddRobot(robot)

		gbot.Start()
	}

For further information refer to keyboard README:
https://github.com/hybridgroup/gobot/blob/master/platforms/keyboard/README.md
*/
package keyboard
