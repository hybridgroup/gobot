// Joystick scanner
// Based on original code from Jacky Boen
// https://github.com/veandco/go-sdl2/blob/master/examples/events/events.go

package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var joysticks [16]*sdl.Joystick

func run() int {
	var event sdl.Event
	var running bool

	sdl.Init(sdl.INIT_JOYSTICK)
	defer sdl.Quit()

	sdl.JoystickEventState(sdl.ENABLE)

	running = true
	for running {
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.JoyAxisEvent:
				fmt.Printf("[%d ms] Axis: %d\tvalue:%d\n",
					t.Timestamp, t.Axis, t.Value)
			case *sdl.JoyBallEvent:
				fmt.Printf("[%d ms] Ball:%d\txrel:%d\tyrel:%d\n",
					t.Timestamp, t.Ball, t.XRel, t.YRel)
			case *sdl.JoyButtonEvent:
				fmt.Printf("[%d ms] Button:%d\tstate:%d\n",
					t.Timestamp, t.Button, t.State)
			case *sdl.JoyHatEvent:
				fmt.Printf("[%d ms] Hat:%d\tvalue:%d\n",
					t.Timestamp, t.Hat, t.Value)
			case *sdl.JoyDeviceEvent:
				if t.Type == sdl.JOYDEVICEADDED {
					joysticks[int(t.Which)] = sdl.JoystickOpen(t.Which)
					if joysticks[int(t.Which)] != nil {
						fmt.Printf("Joystick %d connected\n", t.Which)
					}
				} else if t.Type == sdl.JOYDEVICEREMOVED {
					if joystick := joysticks[int(t.Which)]; joystick != nil {
						joystick.Close()
					}
					fmt.Printf("Joystick %d disconnected\n", t.Which)
				}
			default:
				fmt.Printf("Unknown event\n")
			}
		}

		sdl.Delay(16)
	}

	return 0
}

func main() {
	os.Exit(run())
}
