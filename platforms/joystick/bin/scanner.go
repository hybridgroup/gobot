//go:build utils
// +build utils

// Do not build by default.
//
// Joystick scanner
// Based on original code from
// https://github.com/0xcafed00d/joystick/blob/master/joysticktest/joysticktest.go
// Simple program that displays the state of the specified joystick
//
//	go run joysticktest.go 2
//
// displays state of joystick id 2
package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/0xcafed00d/joystick"
	"github.com/nsf/termbox-go"
)

func printAt(x, y int, s string) {
	for _, r := range s {
		termbox.SetCell(x, y, r, termbox.ColorDefault, termbox.ColorDefault)
		x++
	}
}

func readJoystick(js joystick.Joystick) {
	jinfo, err := js.Read()
	if err != nil {
		printAt(1, 5, "Error: "+err.Error())
		return
	}

	printAt(1, 5, "Buttons:")
	for button := 0; button < js.ButtonCount(); button++ {
		if jinfo.Buttons&(1<<uint32(button)) != 0 {
			printAt(10+button, 5, "X")
			printAt(1, 6, fmt.Sprintf("Button %2d Pressed", button))
		} else {
			printAt(10+button, 5, ".")
		}
	}

	for axis := 0; axis < js.AxisCount(); axis++ {
		printAt(1, axis+8, fmt.Sprintf("Axis %2d Value: %7d", axis, jinfo.AxisData[axis]))
	}

	return
}

func main() {
	jsid := 0
	if len(os.Args) > 1 {
		i, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		jsid = i
	}

	js, jserr := joystick.Open(jsid)

	if jserr != nil {
		fmt.Println(jserr)
		return
	}

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	ticker := time.NewTicker(time.Millisecond * 40)

	for doQuit := false; !doQuit; {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				if ev.Ch == 'q' {
					doQuit = true
				}
			}
			if ev.Type == termbox.EventResize {
				termbox.Flush()
			}

		case <-ticker.C:
			printAt(1, 0, "-- Press 'q' to Exit --")
			printAt(1, 1, fmt.Sprintf("Joystick Name: %s", js.Name()))
			printAt(1, 2, fmt.Sprintf("   Axis Count: %d", js.AxisCount()))
			printAt(1, 3, fmt.Sprintf(" Button Count: %d", js.ButtonCount()))
			readJoystick(js)
			termbox.Flush()
		}
	}
}
