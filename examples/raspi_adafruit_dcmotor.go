package main

import (
	"log"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

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

			dcMotor := 2 // 0-based
			adafruitDCMotorRunner(adaFruit, dcMotor)
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
