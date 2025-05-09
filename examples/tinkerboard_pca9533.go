//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/platforms/asus/tinkerboard"
)

// Wiring
// PWR  Tinkerboard: 1 (+3.3V, VCC), 6, 9, 14, 20 (GND)
// I2C1 Tinkerboard: 3 (SDA), 5 (SCL)
// PCA9533: 8 (VDD, +2.3..5.5V), 4 (VSS, GND), 7 (SDA), 6 (SCL)
// LED pins: 1 (LED0), 2 (LED1), 3 (LED2), 5 (LED3)
// LED's directly driven with pull-up resistors to VDD, e.g. 180 Ohm
// I2C addresses: 0x62 (PCA9533/01), 0x63 (PCA9533/02)
func main() {
	board := tinkerboard.NewAdaptor()
	pca := i2c.NewPCA953xDriver(board, i2c.WithAddress(0x63))

	led := uint8(0)  // index of LED
	wVal := uint8(1) // start with LED is "off"
	rVal := uint8(0)

	work := func() {
		// LED 2 with 2 Hz 1:1, LED 3 with 5 Hz 1:10
		initialize(pca, 2, 5)

		gobot.Every(2000*time.Millisecond, func() {
			fmt.Printf("set LED%d output to %d", led, wVal)
			err := pca.WriteGPIO(led, wVal)
			if err != nil {
				fmt.Println("errW:", err)
			}

			rVal, err = pca.ReadGPIO(led)
			if err != nil {
				fmt.Println("errR:", err)
			}
			if rVal == 0 {
				fmt.Printf(" - LED%d is ON\n", led)
			} else {
				fmt.Printf(" - LED%d is OFF\n", led)
			}

			led = led + 1
			if led > 1 {
				led = 0
				if wVal == 1 {
					wVal = 0
				} else {
					wVal = 1
				}
			}
		})
	}

	robot := gobot.NewRobot("ledI2c",
		[]gobot.Connection{board},
		[]gobot.Device{pca},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}

func initialize(pca *i2c.PCA953xDriver, led2FrequHz float32, led3FrequHz float32) {
	// prepare PWM0
	err := pca.WriteFrequency(0, led2FrequHz)
	if err != nil {
		fmt.Println("errWF0:", err)
	}
	frq, err := pca.ReadFrequency(0)
	if err != nil {
		fmt.Println("errRF0:", err)
	}
	fmt.Println("get Frq0:", frq)

	err = pca.WriteDutyCyclePercent(0, 50)
	if err != nil {
		fmt.Println("errWD0:", err)
	}
	dc, err := pca.ReadDutyCyclePercent(0)
	if err != nil {
		fmt.Println("errRD0:", err)
	}
	fmt.Println("get dc0:", dc)

	// prepare PWM1
	err = pca.WriteFrequency(1, led3FrequHz)
	if err != nil {
		fmt.Println("errWF1:", err)
	}
	frq, err = pca.ReadFrequency(1)
	if err != nil {
		fmt.Println("errRF1:", err)
	}
	fmt.Println("get Frq1:", frq)

	err = pca.WriteDutyCyclePercent(1, 10)
	if err != nil {
		fmt.Println("errWD1:", err)
	}
	dc, err = pca.ReadDutyCyclePercent(1)
	if err != nil {
		fmt.Println("errRD1:", err)
	}
	fmt.Println("get dc1:", dc)

	// LED 2
	fmt.Println("set LED: 2 to: pwm0")
	err = pca.SetLED(2, i2c.PCA953xModePwm0)
	if err != nil {
		fmt.Println("errW:", err)
	}
	fmt.Println("set LED: 3 to: pwm1")
	err = pca.SetLED(3, i2c.PCA953xModePwm1)
	if err != nil {
		fmt.Println("errW:", err)
	}
}
