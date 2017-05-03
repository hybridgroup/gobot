package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
	"log"
)

func main() {

	r := raspi.NewAdaptor()
	ina := i2c.NewINA3221Driver(r)

	work := func() {

		gobot.Every(5*time.Second, func() {
			bv, err := ina.GetBusVoltage(i2c.INA3221Channel1)
			if err != nil {

			}
			log.Printf("Ch 1 Bus Voltage: %f", bv)

			sv, err := ina.GetShuntVoltage(i2c.INA3221Channel1)
			if err != nil {

			}
			log.Printf("Ch 1 Shunt Voltage: %f", sv)

			ma, err := ina.GetCurrent(i2c.INA3221Channel1)
			if err != nil {

			}
			log.Printf("Ch 1 Current: %f", ma)

			lv, err := ina.GetLoadVoltage(i2c.INA3221Channel1)
			if err != nil {

			}
			log.Printf("Ch 1 Load Voltage: %f", lv)
		})
	}

	robot := gobot.NewRobot("ina3221Robot",
		[]gobot.Connection{r},
		[]gobot.Device{ina},
		work,
	)

	robot.Start()
}
