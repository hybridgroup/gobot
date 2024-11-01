//go:build example
// +build example

//
// Do not build by default.

//nolint:gosec // ok here
package main

import (
	"fmt"
	"log"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

const (
	// Min pulse length out of 4096
	servoMin = 150
	// Max pulse length out of 4096
	servoMax = 700
	// Limiting the max this servo can rotate (in deg)
	maxDegree = 180
	// Number of degrees to increase per call
	degIncrease = 10
	yawDeg      = 90
)

func main() {
	r := raspi.NewAdaptor()

	// Changing from the default 0x40 address because this configuration involves
	// a Servo HAT stacked on top of a DC/Stepper Motor HAT on top of the Pi.
	stackedHatAddr := 0x41

	adaFruit := i2c.NewAdafruit2327Driver(r, i2c.WithAddress(stackedHatAddr))

	work := func() {
		gobot.Every(5*time.Second, func() {
			if err := adafruitServoMotorRunner(adaFruit); err != nil {
				fmt.Println(err)
			}
		})
	}

	robot := gobot.NewRobot("adaFruitBot",
		[]gobot.Connection{r},
		[]gobot.Device{adaFruit},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}

func adafruitServoMotorRunner(a *i2c.Adafruit2327Driver) error {
	log.Printf("Servo Motor Run Loop...\n")

	var channel byte = 1
	deg := 90

	// Do not need to set this every run loop
	freq := 60.0
	if err := a.SetServoMotorFreq(freq); err != nil {
		return err
	}
	// start in the middle of the 180-deg range
	pulse := degree2pulse(deg)
	if err := a.SetServoMotorPulse(channel, 0, pulse); err != nil {
		return err
	}
	// INCR
	pulse = degree2pulse(deg + degIncrease)
	if err := a.SetServoMotorPulse(channel, 0, pulse); err != nil {
		return err
	}
	time.Sleep(2000 * time.Millisecond)
	// DECR
	pulse = degree2pulse(deg - degIncrease)
	if err := a.SetServoMotorPulse(channel, 0, pulse); err != nil {
		return err
	}
	return nil
}

func degree2pulse(deg int) int32 {
	pulse := servoMin
	pulse += ((servoMax - servoMin) / maxDegree) * deg
	return int32(pulse)
}
