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
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/asus/tinkerboard"
)

const (
	clockPinOutputID = "24"
	dataPinInputID   = "26"
	calibrationWait  = 10 * time.Second
	cycleTime        = 1 * time.Second
)

// please adjust to your load used for calibration
const (
	setLoad = 100 // used unit will equal the measured unit
	unit    = "gram"
)

// Wiring
// PWR  Tinkerboard: 2(+5V), 6, 9, 14, 20 (GND)
// GPIO Tinkerboard:
// * header pin 24 is the clock output
// * pin 26 is used as data input, an external pull up resistor is suggested (e.g. 10K)
// HX711: the power (VCC) is wired to +5V and GND of tinkerboard, the same for SCK output and the DT input pin
func main() {
	a := tinkerboard.NewAdaptor()
	hx711 := gpio.NewHX711Driver(a, clockPinOutputID, dataPinInputID,
		gpio.WithHX711ReadTimeout(1*time.Second), gpio.WithHX711ReadTriesMax(10))

	work := func() {
		fmt.Println("First lets calibrate the sensor...")
		fmt.Printf("Please remove the mass completely for zeroing within %s\n", calibrationWait)
		time.Sleep(calibrationWait)

		fmt.Println("Zeroing starts")
		if err := hx711.Zero(false); err != nil {
			fmt.Println("Zeroing failed")
			panic(err)
		}

		fmt.Printf("Please apply the load (%d%s) for calibration within %s\n", setLoad, unit, calibrationWait)
		time.Sleep(calibrationWait)

		fmt.Println("Calibration starts")
		if err := hx711.Calibrate(setLoad, false); err != nil {
			fmt.Println("Calibration failed")
			panic(err)
		}

		offs, factor := hx711.OffsetAndCalibrationFactor(false)
		fmt.Printf("Calibration done completely, offset: %d factor: %f\n", offs, factor)

		fmt.Println("Now start repeated measurement...")
		_ = gobot.Every(cycleTime, func() {
			if v, _, err := hx711.Measure(); err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("measure done: %5.3f gram\n", v)
			}
		})
	}

	robot := gobot.NewRobot("loadCellBot",
		[]gobot.Connection{a},
		[]gobot.Device{hx711},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
