// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/digispark"
)

func main() {
	board := digispark.NewAdaptor()
	mpl115a2 := i2c.NewMPL115A2Driver(board)

	work := func() {
		gobot.Every(1*time.Second, func() {
			press, _ := mpl115a2.Pressure()
			fmt.Println("Pressure", press)

			temp, _ := mpl115a2.Temperature()
			fmt.Println("Temperature", temp)
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{board},
		[]gobot.Device{mpl115a2},
		work,
	)

	err := robot.Start()
	if err != nil {
		fmt.Println(err)
	}
}
