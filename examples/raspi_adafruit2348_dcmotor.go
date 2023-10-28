//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"log"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()

	// we use the default address 0x60 for DC/Stepper Motor HAT on top of the Pi
	adaFruit := i2c.NewAdafruit2348Driver(r)

	work := func() {
		gobot.Every(5*time.Second, func() {
			dcMotor := 2 // 0-based
			adafruitDCMotorRunner(adaFruit, dcMotor)
		})
	}

	robot := gobot.NewRobot("adaFruitBot",
		[]gobot.Connection{r},
		[]gobot.Device{adaFruit},
		work,
	)

	robot.Start()
}

func adafruitDCMotorRunner(a *i2c.Adafruit2348Driver, dcMotor int) (err error) {
	log.Printf("DC Motor Run Loop...\n")
	// set the speed:
	var speed int32 = 255 // 255 = full speed!
	if err = a.SetDCMotorSpeed(dcMotor, speed); err != nil {
		return
	}
	// run FORWARD
	if err = a.RunDCMotor(dcMotor, i2c.Adafruit2348Forward); err != nil {
		return
	}
	// Sleep and RELEASE
	time.Sleep(2000 * time.Millisecond)
	if err = a.RunDCMotor(dcMotor, i2c.Adafruit2348Release); err != nil {
		return
	}
	// run BACKWARD
	if err = a.RunDCMotor(dcMotor, i2c.Adafruit2348Backward); err != nil {
		return
	}
	// Sleep and RELEASE
	time.Sleep(2000 * time.Millisecond)
	if err = a.RunDCMotor(dcMotor, i2c.Adafruit2348Release); err != nil {
		return
	}
	return
}
