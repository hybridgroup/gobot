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

func CCS811BootData(a *i2c.CCS811Driver) {
	v, err := a.GetHardwareVersion()
	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Printf("Hardare Version %#x\n", v)

	d, err := a.GetFirmwareBootVersion()
	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Printf("Boot Version %#x\n", d)

	d, err = a.GetFirmwareAppVersion()
	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Printf("App Version %#x\n\n", d)
}

func main() {
	r := raspi.NewAdaptor()
	ccs811Driver := i2c.NewCCS811Driver(r)

	work := func() {
		CCS811BootData(ccs811Driver)
		gobot.Every(1*time.Second, func() {
			s, err := ccs811Driver.GetStatus()
			if err != nil {
				fmt.Printf("Error fetching data from the status register: %+v\n", err.Error())
			}
			fmt.Printf("Status %+v \n", s)

			hd, err := ccs811Driver.HasData()
			if err != nil {
				fmt.Printf("Error fetching data from the status register: %+v\n", err.Error())
			}

			if hd {
				ec02, tv0C, _ := ccs811Driver.GetGasData()
				fmt.Printf("Gas Data - eco2: %+v, tvoc: %+v \n", ec02, tv0C)

				temp, _ := ccs811Driver.GetTemperature()
				fmt.Printf("Temperature %+v \n\n", temp)
			} else {
				fmt.Println("New data is not avaliable\n")
			}
		})
	}

	robot := gobot.NewRobot("adaFruitBot",
		[]gobot.Connection{r},
		[]gobot.Device{ccs811Driver},
		work,
	)

	robot.Start()
}
