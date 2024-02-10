//go:build example
// +build example

//
// Do not build by default.

/*
How to run:
Connect to the drone's Wi-Fi network from your computer. It will be named something like "TELLO-XXXXXX".

Once you are connected you can run the Gobot code on your computer to control the drone.

	go run examples/tello_keyboard.go
*/

package main

import (
	"fmt"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/dji/tello"
	"gobot.io/x/gobot/v2/platforms/keyboard"
)

func resetDronePostion(drone *tello.Driver) {
	if err := drone.Forward(0); err != nil {
		fmt.Println(err)
	}
	if err := drone.Backward(0); err != nil {
		fmt.Println(err)
	}
	if err := drone.Up(0); err != nil {
		fmt.Println(err)
	}
	if err := drone.Down(0); err != nil {
		fmt.Println(err)
	}
	if err := drone.Left(0); err != nil {
		fmt.Println(err)
	}
	if err := drone.Right(0); err != nil {
		fmt.Println(err)
	}
	if err := drone.Clockwise(0); err != nil {
		fmt.Println(err)
	}
}

func main() {
	drone := tello.NewDriver("8888")
	keys := keyboard.NewDriver()

	_ = keys.On(keyboard.Key, func(data interface{}) {
		key := data.(keyboard.KeyEvent)
		switch key.Key {
		case keyboard.A:
			fmt.Println(key.Char)
			if err := drone.Clockwise(-25); err != nil {
				fmt.Println(err)
			}
		case keyboard.D:
			fmt.Println(key.Char)
			if err := drone.Clockwise(25); err != nil {
				fmt.Println(err)
			}
		case keyboard.W:
			fmt.Println(key.Char)
			if err := drone.Forward(20); err != nil {
				fmt.Println(err)
			}
		case keyboard.S:
			fmt.Println(key.Char)
			if err := drone.Backward(20); err != nil {
				fmt.Println(err)
			}
		case keyboard.K:
			fmt.Println(key.Char)
			if err := drone.Down(20); err != nil {
				fmt.Println(err)
			}
		case keyboard.J:
			fmt.Println(key.Char)
			if err := drone.Up(20); err != nil {
				fmt.Println(err)
			}
		case keyboard.Q:
			fmt.Println(key.Char)
			if err := drone.Land(); err != nil {
				fmt.Println(err)
			}
		case keyboard.P:
			fmt.Println(key.Char)
			if err := drone.TakeOff(); err != nil {
				fmt.Println(err)
			}
		case keyboard.ArrowUp:
			fmt.Println(key.Char)
			if err := drone.FrontFlip(); err != nil {
				fmt.Println(err)
			}
		case keyboard.ArrowDown:
			fmt.Println(key.Char)
			if err := drone.BackFlip(); err != nil {
				fmt.Println(err)
			}
		case keyboard.ArrowLeft:
			fmt.Println(key.Char)
			if err := drone.LeftFlip(); err != nil {
				fmt.Println(err)
			}
		case keyboard.ArrowRight:
			fmt.Println(key.Char)
			if err := drone.RightFlip(); err != nil {
				fmt.Println(err)
			}
		case keyboard.Escape:
			resetDronePostion(drone)
		}
	})

	var flightData *tello.FlightData
	work := func() {
		_ = drone.On(tello.FlightDataEvent, func(data interface{}) {
			flightData = data.(*tello.FlightData)
			fmt.Println("Height:", flightData.Height)
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{keys, drone},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
