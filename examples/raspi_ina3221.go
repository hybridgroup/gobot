// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {

	r := raspi.NewAdaptor()
	ina := i2c.NewINA3221Driver(r)

	work := func() {

		gobot.Every(5*time.Second, func() {
			for _, ch := range []i2c.INA3221Channel{i2c.INA3221Channel1, i2c.INA3221Channel2, i2c.INA3221Channel3} {
				val, err := ina.GetBusVoltage(ch)
				if err != nil {
					fmt.Printf("INA3221Channel %v Bus Voltage error: %v\n", ch, err)
				}
				fmt.Printf("INA3221Channel %v Bus Voltage: %fV\n", ch, val)

				val, err = ina.GetShuntVoltage(ch)
				if err != nil {
					fmt.Printf("INA3221Channel %v Shunt Voltage error: %v\n", ch, err)
				}
				fmt.Printf("INA3221Channel %v Shunt Voltage: %fV\n", ch, val)

				val, err = ina.GetCurrent(ch)
				if err != nil {
					fmt.Printf("INA3221Channel %v Current error: %v\n", ch, err)
				}
				fmt.Printf("INA3221Channel %v Current: %fmA\n", ch, val)

				val, err = ina.GetLoadVoltage(ch)
				if err != nil {
					fmt.Printf("INA3221Channel %v Load Voltage error: %v\n", ch, err)
				}
				fmt.Printf("INA3221Channel %v Load Voltage: %fV\n", ch, val)
			}
		})
	}

	robot := gobot.NewRobot("INA3221 Robot",
		[]gobot.Connection{r},
		[]gobot.Device{ina},
		work,
	)

	robot.Start()
}
