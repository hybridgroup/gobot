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
	a.SetStepperMotorSpeed(motor, speed)

	//fmt.Println("Single coil steps")
	if err = a.Step(motor, 100, i2c.AdafruitForward, i2c.AdafruitSingle); err != nil {
		log.Printf(err.Error())
		return
	}
	if err = a.Step(motor, 100, i2c.AdafruitBackward, i2c.AdafruitSingle); err != nil {
		log.Printf(err.Error())
		return
	}
	return
}
func adafruitDCMotorRunner(a *i2c.AdafruitMotorHatDriver, dcMotor int) (err error) {

	log.Printf("DC Motor Run Loop...\n")
	// set the speed:
	var speed int32 = 255 // 255 = full speed!
	if err = a.SetDCMotorSpeed(dcMotor, speed); err != nil {
		return
	}
	// run FORWARD
	if err = a.RunDCMotor(dcMotor, i2c.AdafruitForward); err != nil {
		return
	}
	// Sleep and RELEASE
	<-time.After(2000 * time.Millisecond)
	if err = a.RunDCMotor(dcMotor, i2c.AdafruitRelease); err != nil {
		return
	}
	// run BACKWARD
	if err = a.RunDCMotor(dcMotor, i2c.AdafruitBackward); err != nil {
		return
	}
	// Sleep and RELEASE
	<-time.After(2000 * time.Millisecond)
	if err = a.RunDCMotor(dcMotor, i2c.AdafruitRelease); err != nil {
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

			//dcMotor := 3 // 0-based
			//adafruitDCMotorRunner(adaFruit, dcMotor)

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
