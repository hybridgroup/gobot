// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"gobot.io/x/gobot/sysfs"
)

func main() {
	r := raspi.NewAdaptor()

	motor := gpio.NewMotorDriver(r, "11")
	motor.ForwardPin = "40"
	motor.BackwardPin = "38"

	err := motor.Backward(255)
	if err != nil {
		panic(err)
	}

	defer func() {
		motor.Off()
	}()

	encoder, err := r.DigitalPin("36", "")
	if err != nil {
		panic(err)
	}

	listener, err := sysfs.NewInterruptListener()
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	listener.Start()

	numEvents := 0

	err = encoder.Listen("falling", listener, func(b byte) {
		numEvents++
	})
	if err != nil {
		panic(err)
	}

	secs := time.Second * 2
	time.Sleep(secs)

	fmt.Printf("numEvents: %d\n", numEvents)
}
