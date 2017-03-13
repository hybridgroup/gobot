// +build example
//
// Do not build by default.

package main

import (
	"log"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func adafruitStepperMotorRunner(a *i2c.AdafruitMotorHatDriver, motor int) (err error) {
	log.Printf("Stepper Motor Run Loop...\n")
	// set the speed state:
	speed := 30 // rpm
	style := i2c.AdafruitDouble
	steps := 20

	a.SetStepperMotorSpeed(motor, speed)

	if err = a.Step(motor, steps, i2c.AdafruitForward, style); err != nil {
		log.Printf(err.Error())
		return
	}
	if err = a.Step(motor, steps, i2c.AdafruitBackward, style); err != nil {
		log.Printf(err.Error())
		return
	}
	return
}

func main() {
	r := raspi.NewAdaptor()
	adaFruit := i2c.NewAdafruitMotorHatDriver(r)

	work := func() {
		gobot.Every(5*time.Second, func() {
			motor := 0 // 0-based
			adafruitStepperMotorRunner(adaFruit, motor)
		})
	}

	robot := gobot.NewRobot("adaFruitBot",
		[]gobot.Connection{r},
		[]gobot.Device{adaFruit},
		work,
	)

	robot.Start()
}
