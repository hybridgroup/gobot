//go:build example
// +build example

//
// Do not build by default.

//nolint:gosec // ok here
package main

import (
	"log"
	"os"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

func main() {
	const (
		coilA1 = "7"
		coilA2 = "13"
		coilB1 = "11"
		coilB2 = "15"

		degPerStep = 1.875
		countRot   = 10
	)
	stepPerRevision := int(360.0 / degPerStep)

	r := raspi.NewAdaptor()
	stepper := gpio.NewStepperDriver(r, [4]string{coilA1, coilB1, coilA2, coilB2}, gpio.StepperModes.DualPhaseStepping,
		uint(stepPerRevision))

	work := func() {
		defer func() {
			ec := 0
			// set current to zero to prevent overheating
			if err := stepper.Sleep(); err != nil {
				ec = 1
				log.Println("work done", err)
			} else {
				log.Println("work done")
			}

			os.Exit(ec)
		}()

		gobot.After(5*time.Second, func() {
			// this stops only the current movement and the next will start immediately (if any)
			// this means for the example, that the first rotation stops after ~5 rotations
			log.Println("asynchron stop after 5 sec.")
			if err := stepper.Stop(); err != nil {
				log.Println(err)
			}
		})

		// one rotation per second
		if err := stepper.SetSpeed(60); err != nil {
			log.Println("set speed", err)
		}

		// Move forward N revolution
		if err := stepper.Move(stepPerRevision * countRot); err != nil {
			log.Println("move forward", err)
		}

		// Move backward N revolution
		if err := stepper.MoveDeg(-360 * countRot); err != nil {
			log.Println("move backward", err)
		}
	}

	robot := gobot.NewRobot("stepperBot",
		[]gobot.Connection{r},
		[]gobot.Device{stepper},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
