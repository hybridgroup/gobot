package main

import (
	"log"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/hybridgroup/gobot/platforms/raspi"
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
	gbot := gobot.NewGobot()
	r := raspi.NewRaspiAdaptor("raspi")
	adaFruit := i2c.NewAdafruitMotorHatDriver(r, "adafruit")

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

	gbot.AddRobot(robot)

	gbot.Start()
}
