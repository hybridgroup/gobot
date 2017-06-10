# Keyboard

This module implements support for keyboard events, wrapping the `stty` utility.

## How to Install

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use

Example parsing key events

```go
package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/keyboard"
)

func main() {
	keys := keyboard.NewDriver()

	work := func() {
		keys.On(keyboard.Key, func(data interface{}) {
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

	robot.Start()
}
```
