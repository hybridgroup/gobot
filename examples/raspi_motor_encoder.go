// +build example
//
// Do not build by default.

package main

import (
	"log"
	"time"

	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"gobot.io/x/gobot/sysfs"
)

func main() {
	err := mainReal()
	if err != nil {
		log.Fatal(err)
	}
}

func mainReal() error {
	r := raspi.NewAdaptor()

	motor := gpio.NewMotorDriver(r, "11")
	motor.ForwardPin = "40"
	motor.BackwardPin = "38"

	defer func() {
		motor.Off()
	}()

	err := motor.Backward(255)
	if err != nil {
		return err
	}

	r.DigitalPinSetPullUpDown("36", true)
	encoder, err := r.DigitalPin("36", "")
	if err != nil {
		return err
	}

	listener, err := sysfs.NewInterruptListener()
	if err != nil {
		return err
	}
	defer listener.Close()
	listener.Start()

	numEvents := 0

	err = encoder.Listen("falling", listener, func(b byte) {
		numEvents++
	})
	if err != nil {
		return err
	}

	secs := time.Second * 2
	time.Sleep(secs)

	log.Printf("numEvents: %d", numEvents)
	return nil
}
