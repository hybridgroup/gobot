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

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gobot.io/x/gobot/platforms/keyboard"
)

func resetDronePostion(drone *tello.Driver) {
	drone.Forward(0)
	drone.Backward(0)
	drone.Up(0)
	drone.Down(0)
	drone.Left(0)
	drone.Right(0)
	drone.Clockwise(0)
}

func main() {
	drone := tello.NewDriver("8888")
	keys := keyboard.NewDriver()

	keys.On(keyboard.Key, func(data interface{}) {
		key := data.(keyboard.KeyEvent)
		switch key.Key {
		case keyboard.A:
			fmt.Println(key.Char)
			drone.Clockwise(-25)
		case keyboard.D:
			fmt.Println(key.Char)
			drone.Clockwise(25)
		case keyboard.W:
			fmt.Println(key.Char)
			drone.Forward(20)
		case keyboard.S:
			fmt.Println(key.Char)
			drone.Backward(20)
		case keyboard.K:
			fmt.Println(key.Char)
			drone.Down(20)
		case keyboard.J:
			fmt.Println(key.Char)
			drone.Up(20)
		case keyboard.Q:
			fmt.Println(key.Char)
			drone.Land()
		case keyboard.P:
			fmt.Println(key.Char)
			drone.TakeOff()
		case keyboard.ArrowUp:
			fmt.Println(key.Char)
			drone.FrontFlip()
		case keyboard.ArrowDown:
			fmt.Println(key.Char)
			drone.BackFlip()
		case keyboard.ArrowLeft:
			fmt.Println(key.Char)
			drone.LeftFlip()
		case keyboard.ArrowRight:
			fmt.Println(key.Char)
			drone.RightFlip()
		case keyboard.Escape:
			resetDronePostion(drone)
		}
	})

	var flightData *tello.FlightData
	work := func() {
		drone.On(tello.FlightDataEvent, func(data interface{}) {
			flightData = data.(*tello.FlightData)
			fmt.Println("Height:", flightData.Height)
		})
	}

	robot := gobot.NewRobot("tello",
		[]gobot.Connection{},
		[]gobot.Device{keys, drone},
		work,
	)

	robot.Start()
}
