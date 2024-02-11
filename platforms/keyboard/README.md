# Keyboard

This module implements support for keyboard events, wrapping the `stty` utility.

## How to Install

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

## How to Use

Example parsing key events

```go
package main

import (
  "fmt"

  "gobot.io/x/gobot/v2"
  "gobot.io/x/gobot/v2/platforms/keyboard"
)

func main() {
  keys := keyboard.NewDriver()

  work := func() {
    _ = keys.On(keyboard.Key, func(data interface{}) {
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

  if err := robot.Start(); err != nil {
		panic(err)
	}
}
```
