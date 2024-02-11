//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
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
			motor := 0 // 0-based
			if err := adafruitStepperMotorRunner(adaFruit, motor); err != nil {
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

func adafruitStepperMotorRunner(a *i2c.Adafruit2348Driver, motor int) error {
	log.Printf("Stepper Motor Run Loop...\n")
	// set the speed state:
	speed := 30 // rpm
	style := i2c.Adafruit2348Double
	steps := 20

	if err := a.SetStepperMotorSpeed(motor, speed); err != nil {
		return err
	}

	if err := a.Step(motor, steps, i2c.Adafruit2348Forward, style); err != nil {
		return err
	}

	if err := a.Step(motor, steps, i2c.Adafruit2348Backward, style); err != nil {
		return err
	}

	return nil
}
