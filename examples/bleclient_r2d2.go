//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass the Bluetooth address or name as the first param:

	go run examples/bleclient_r2d2.go D2-1234

 NOTE: sudo is required to use BLE in Linux
 NOTE: GODEBUG=cgocheck=0 is required in the prefix of the command to use BLE in OsX
*/

//nolint:gosec // ok here
package main

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/sphero"
	"gobot.io/x/gobot/v2/platforms/bleclient"
	"log"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.Lmicroseconds)

	bleAdaptor := bleclient.NewAdaptor(os.Args[1])
	r2 := sphero.NewR2D2Driver(bleAdaptor)

	work := func() {
		// this library is not opinionated about initializing the R2, perform these steps if needed
		r2.ResetLocatorData()
		r2.SetDomePosition(0)

		// you can use gobot.Rand(255) to get a random color value up to 255
		// or set the values yourself
		//r := uint8(180)
		//g := uint8(180)
		//b := uint8(45)
		//r2.SetRGB(r, g, b)
		//time.Sleep(1 * time.Second)
		//r2.SetDomePosition(160)
		//time.Sleep(1 * time.Second)
		//r2.PlaySound(sphero.R2Alarm4, sphero.Immediate)
		//time.Sleep(500 * time.Millisecond)
		//r2.SetDomePosition(-90)

		// make sure you perform a ThreeLegs action before trying to Roll
		//r2.PerformLegAction(sphero.ThreeLegs)
		//time.Sleep(4 * time.Second)
		//r2.Roll(100, 0)
		//time.Sleep(2 * time.Second)
		//r2.Stop()
		//time.Sleep(1 * time.Second)
		//r2.PerformLegAction(sphero.TwoLegs)

		// use the AfterCurrentSound to play sounds one after another
		r2.PlaySound(sphero.R2Chatty4, sphero.AfterCurrentSound)
		r2.PlaySound(sphero.R2Positive7, sphero.AfterCurrentSound)
		// wait for sounds to play before using an Animation
		time.Sleep(3 * time.Second)
		r2.PlayAnimation(sphero.EmoteExcited)
	}

	robot := gobot.NewRobot("r2Bot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{r2},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
